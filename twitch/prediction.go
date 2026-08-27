package twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/monktype/msc/keys"
	"github.com/nicklaw5/helix/v2"
)

// CreatePrediction creates a Twitch prediction with the given title, window, and choices.
// Returns a prediction ID and error.
// Twitch predictions support exactly 2 outcomes (each max 25 characters); the title is
// max 45 characters and the prediction window is 1..1800 seconds.
func CreatePrediction(c *helix.Client, channelID string, title string, windowInSeconds int, choices []string) (string, error) {
	// Convert choices to a slice of PredictionChoiceParam
	var predictionChoices []helix.PredictionChoiceParam
	for _, choice := range choices {
		if len(choice) > 25 {
			fmt.Printf("Warning: Choice '%s' exceeds the maximum length of 25 characters.\n", choice)
			choice = choice[:25] // Truncate if necessary
		}
		predictionChoices = append(predictionChoices, helix.PredictionChoiceParam{Title: choice})
	}

	prediction, err := c.CreatePrediction(&helix.CreatePredictionParams{
		BroadcasterID:    channelID,
		Title:            title,
		Outcomes:         predictionChoices,
		PredictionWindow: windowInSeconds,
	})
	if err != nil {
		fmt.Printf("Creating a prediction failed: %s\n", err)
		return "", err
	}
	if prediction.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", prediction)
		return "", fmt.Errorf("check status code information")
	}

	fmt.Printf("Prediction created with ID: %s\n", prediction.Data.Predictions[0].ID)
	return prediction.Data.Predictions[0].ID, nil
}

// PredictionsPage is one page of predictions returned by GetPredictions.
// Cursor is the cursor to pass as `after` to GetPredictions to fetch the next page;
// it is "" when there is no next page.
type PredictionsPage struct {
	Predictions []helix.Prediction `json:"predictions"`
	Cursor      string             `json:"cursor"`
}

// GetPredictions gets predictions from a channel ID.
// `after` is an optional pagination cursor to continue from a previous page; the
// returned page's Cursor carries the cursor for the next page ("" when done).
func GetPredictions(c *helix.Client, channelID string, after string) (PredictionsPage, error) {
	var page PredictionsPage

	predictions, err := c.GetPredictions(&helix.PredictionsParams{
		BroadcasterID: channelID,
		After:         after,
	})
	if err != nil {
		fmt.Printf("Failed to get predictions on channel ID %s: %s\n", channelID, err)
		return page, err
	}
	if predictions.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", predictions)
		return page, fmt.Errorf("check status code information")
	}

	page.Predictions = predictions.Data.Predictions
	page.Cursor = predictions.Data.Pagination.Cursor
	return page, nil
}

// GetAllPredictions returns every prediction on the channel by fetching all pages.
// It is the convenience wrapper around GetPredictions for callers that don't want
// to manage the cursor loop themselves. It keeps going until the cursor is
// exhausted; if the cursor ever stops advancing (the API returns a cursor it
// already gave), it stops and returns the accumulated list along with a warning
// rather than looping forever on the same page.
func GetAllPredictions(c *helix.Client, channelID string) ([]helix.Prediction, error) {
	var all []helix.Prediction
	after := ""
	seen := map[string]bool{} // cursors already received; a repeat means we're looping
	for {
		p, err := GetPredictions(c, channelID, after)
		if err != nil {
			return nil, err
		}
		all = append(all, p.Predictions...)
		if p.Cursor == "" {
			break // no next page — done
		}
		if seen[p.Cursor] {
			fmt.Printf("Warning: GetAllPredictions stopped early; the cursor stopped advancing and the list may be incomplete.\n")
			break
		}
		seen[p.Cursor] = true
		after = p.Cursor
	}
	return all, nil
}

// GetPrediction gets a single prediction from a prediction ID string and a channel ID string.
func GetPrediction(c *helix.Client, channelID string, predictionID string) (helix.Prediction, error) {
	var emptyPredictionResponse helix.Prediction

	predictions, err := c.GetPredictions(&helix.PredictionsParams{
		BroadcasterID: channelID,
		ID:            predictionID,
	})
	if err != nil {
		fmt.Printf("Failed to get prediction on channel ID %s: %s\n", channelID, err)
		return emptyPredictionResponse, err
	}
	if predictions.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", predictions)
		return emptyPredictionResponse, fmt.Errorf("check status code information")
	}
	if len(predictions.Data.Predictions) == 0 {
		fmt.Printf("Prediction %s was not found on channel %s.\n", predictionID, channelID)
		return emptyPredictionResponse, nil
	}

	return predictions.Data.Predictions[0], nil
}

// EndPrediction resolves a prediction, awarding channel points to the predictors of the
// winning outcome. `status` is "RESOLVED" and `winningOutcomeID` is the ID of one of the
// prediction's outcomes. Canceling (refunding) is not done here — see CancelPrediction,
// which needs the winning_outcome_id field omitted entirely.
func EndPrediction(c *helix.Client, channelID string, predictionID string, status string, winningOutcomeID string) error {
	resp, err := c.EndPrediction(&helix.EndPredictionParams{
		BroadcasterID:    channelID,
		ID:               predictionID,
		Status:           status,
		WinningOutcomeID: winningOutcomeID,
	})
	if err != nil {
		fmt.Printf("Failed to end prediction: %s\n", err)
		return err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return fmt.Errorf("check status code information")
	}

	return nil
}

// CancelPrediction cancels a prediction with no winner, refunding channel points
// (status CANCELED). See endPredictionNoOutcome for why this isn't a helix call.
func CancelPrediction(channelID string, predictionID string) error {
	return endPredictionNoOutcome(channelID, predictionID, "CANCELED")
}

// LockPrediction moves a running prediction from ACTIVE to LOCKED early, closing the
// prediction window so no further predictions can be placed, without picking a winner
// or refunding (status LOCKED). Once locked, the prediction can be resolved.
// See endPredictionNoOutcome for why this isn't a helix call.
func LockPrediction(channelID string, predictionID string) error {
	return endPredictionNoOutcome(channelID, predictionID, "LOCKED")
}

// endPredictionNoOutcome issues the end-prediction PATCH for the updates that carry no
// winning outcome — CANCELED (refund) and LOCKED (stop new entries).
//
// The endpoint accepts status RESOLVED, CANCELED, or LOCKED; TERMINATED — which the
// nicklaw5/helix doc comment suggests — is poll vocabulary and is rejected here with 400.
//
// It is a raw request rather than helix's EndPrediction because that helper's
// EndPredictionParams has no omitempty on winning_outcome_id, so a no-winner update
// would serialize `"winning_outcome_id": ""`, a field only valid for a RESOLVED update.
// These updates are just broadcaster_id + id + status.
//
// Credentials are read from the keystore (the same source GetClient uses and keeps
// current), so call GetClient first to ensure the token is valid/refreshed.
func endPredictionNoOutcome(channelID string, predictionID string, status string) error {
	clientID, err := keys.GetKey("client-id")
	if err != nil {
		fmt.Printf("Failed to get Client ID from keystore: %s\n", err)
		return err
	}

	accessToken, err := keys.GetKey("access-token")
	if err != nil {
		fmt.Printf("Failed to get access token from keystore: %s\n", err)
		return err
	}

	payload, err := json.Marshal(map[string]string{
		"broadcaster_id": channelID,
		"id":             predictionID,
		"status":         status,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, helix.DefaultAPIBaseURL+"/predictions", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Failed to end prediction: %s\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		var apiErr struct {
			Error   string `json:"error"`
			Status  int    `json:"status"`
			Message string `json:"message"`
		}
		detail := ""
		if json.Unmarshal(body, &apiErr) == nil {
			if apiErr.Message != "" {
				detail = apiErr.Message
			} else {
				detail = apiErr.Error
			}
		}
		if detail == "" {
			detail = string(body)
		}
		fmt.Printf("Status code was bad: %d %s\n", resp.StatusCode, detail)
		return fmt.Errorf("check status code information")
	}

	return nil
}
