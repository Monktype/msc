package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/monktype/msc/twitch"
	"github.com/nicklaw5/helix/v2"
	"github.com/spf13/cobra"
)

var predictionCmd = &cobra.Command{
	Use:   "prediction",
	Short: "Create a prediction with -c (channel name), -d (window in seconds), -t (title), followed by exactly 2 choices. Close it out later with `prediction resolve` or `prediction cancel`, or `prediction lock` to stop new entries early.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			fmt.Printf("Exactly 2 prediction choices are required, %d provided.\n", len(args))
			return fmt.Errorf("give the right number of arguments")
		}

		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		title, err := cmd.Flags().GetString("title")
		if err != nil {
			return err
		}

		window, err := cmd.Flags().GetInt("window")
		if err != nil {
			return err
		}

		if window < 1 || window > 1800 {
			fmt.Printf("Prediction window must be between 1 and 1800 seconds, got %d.\n", window)
			return fmt.Errorf("invalid window duration")
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		userID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		predictionID, err := twitch.CreatePrediction(c, userID, title, window, args)
		if err != nil {
			return err
		}

		// Predictions do not resolve themselves, so print the outcomes and the exact
		// commands to close this one out later (resolve picks a winner; cancel refunds).
		fmt.Printf("Outcomes:\n")
		for i, choice := range args {
			fmt.Printf("  %d) %s\n", i+1, choice)
		}
		fmt.Printf("\n")
		fmt.Printf("Resolve with: msc prediction resolve -c %s -i %s --outcome <1|2|title>\n", channelname, predictionID)
		fmt.Printf("Cancel with:  msc prediction cancel -c %s -i %s\n", channelname, predictionID)
		fmt.Printf("Lock early:   msc prediction lock -c %s -i %s (closes the window before the timer runs out)\n", channelname, predictionID)

		return nil
	},
}

var predictionResolveCmd = &cobra.Command{
	Use:   "resolve",
	Short: "Resolve a prediction by picking the winning outcome (awards channel points to its predictors)",
	RunE: func(cmd *cobra.Command, args []string) error {
		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		predictionID, err := cmd.Flags().GetString("prediction-id")
		if err != nil {
			return err
		}

		outcome, err := cmd.Flags().GetString("outcome")
		if err != nil {
			return err
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		prediction, err := twitch.GetPrediction(c, channelID, predictionID)
		if err != nil {
			return err
		}
		if prediction.ID == "" {
			return fmt.Errorf("prediction %s not found on channel %s", predictionID, channelname)
		}

		winnerIndex := findOutcomeIndex(prediction.Outcomes, outcome)
		if winnerIndex < 0 {
			return fmt.Errorf("outcome %q not found; use 1..%d or an outcome title", outcome, len(prediction.Outcomes))
		}
		winner := prediction.Outcomes[winnerIndex]

		if err := twitch.EndPrediction(c, channelID, predictionID, "RESOLVED", winner.ID); err != nil {
			return err
		}

		fmt.Printf("Prediction %s resolved; outcome %q wins.\n", predictionID, winner.Title)
		return nil
	},
}

var predictionCancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel a prediction with no winner (channel points are refunded)",
	RunE: func(cmd *cobra.Command, args []string) error {
		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		predictionID, err := cmd.Flags().GetString("prediction-id")
		if err != nil {
			return err
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		prediction, err := twitch.GetPrediction(c, channelID, predictionID)
		if err != nil {
			return err
		}
		if prediction.ID == "" {
			return fmt.Errorf("prediction %s not found on channel %s", predictionID, channelname)
		}
		fmt.Printf("Prediction %s is currently in status: %s\n", predictionID, prediction.Status)
		if prediction.Status == "RESOLVED" || prediction.Status == "CANCELED" {
			return fmt.Errorf("prediction %s is already %s; nothing to cancel", predictionID, prediction.Status)
		}

		if err := twitch.CancelPrediction(channelID, predictionID); err != nil {
			return err
		}

		fmt.Printf("Prediction %s canceled; channel points were refunded.\n", predictionID)
		return nil
	},
}

var predictionLockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock a running prediction early (closes the window so no more predictions can be placed; no winner, no refund)",
	RunE: func(cmd *cobra.Command, args []string) error {
		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		predictionID, err := cmd.Flags().GetString("prediction-id")
		if err != nil {
			return err
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		prediction, err := twitch.GetPrediction(c, channelID, predictionID)
		if err != nil {
			return err
		}
		if prediction.ID == "" {
			return fmt.Errorf("prediction %s not found on channel %s", predictionID, channelname)
		}
		fmt.Printf("Prediction %s is currently in status: %s\n", predictionID, prediction.Status)
		if prediction.Status != "ACTIVE" {
			return fmt.Errorf("prediction %s is %s; only a running (ACTIVE) prediction can be locked early", predictionID, prediction.Status)
		}

		if err := twitch.LockPrediction(channelID, predictionID); err != nil {
			return err
		}

		fmt.Printf("Prediction %s locked; the window is closed and no further predictions will be accepted.\n", predictionID)
		return nil
	},
}

// findOutcomeIndex finds the index of the outcome matching `spec`. A numeric value
// ("1", "2") matches an outcome by position; anything else is matched against the
// outcome titles (case-insensitive). Returns -1 if nothing matches.
func findOutcomeIndex(outcomes []helix.Outcomes, spec string) int {
	if n, err := strconv.Atoi(spec); err == nil && n >= 1 && n <= len(outcomes) {
		return n - 1
	}
	for i := range outcomes {
		if strings.EqualFold(outcomes[i].Title, spec) {
			return i
		}
	}
	return -1
}
