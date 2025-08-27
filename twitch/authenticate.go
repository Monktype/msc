package twitch

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/monktype/msc/callback"
	"github.com/monktype/msc/keys"
	"github.com/nicklaw5/helix/v2"
)

type AuthType int

const (
	AuthToken AuthType = 0 // Implicit Grant
	AuthCode  AuthType = 1 // Authorization Code Grant
)

// generateRandomState generates a random URL-safe base64 encoded string.
func generateRandomState(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// Authenticate is the process to get a user token from Twitch.
// It requires that the setup process (putting client ID into the keystore) has already taken place.
func Authenticate(authType AuthType) error {
	// Generate a random state
	state, err := generateRandomState(16)
	if err != nil {
		fmt.Printf("Failed to generate random state: %s\n", err)
		return err
	}

	responseChan := make(chan callback.CallbackResponse)
	shutdownChan := make(chan bool)

	go func() {
		err := callback.StartCallbackServer(state, responseChan, shutdownChan)
		if err != nil {
			log.Fatalf("Failed to start callback server: %s\n", err)
		}
	}()

	clientID, err := keys.GetKey("client-id")
	if err != nil {
		fmt.Printf("Failed to get Client ID from keystore: %s\n", err)
		return err
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:    clientID,
		RedirectURI: "http://localhost:3024/redirect",
	})
	if err != nil {
		fmt.Printf("Unable to create client for authentication: %s\n", err)
		return err
	}

	var authTypeString string

	switch authType {
	case AuthToken:
		authTypeString = "token"
	case AuthCode:
		authTypeString = "code"
	}

	url := client.GetAuthorizationURL(&helix.AuthorizationURLParams{
		ResponseType: authTypeString,
		Scopes: []string{
			"channel:edit:commercial",
			"channel:manage:polls",
			"channel:manage:predictions",
			"channel:manage:redemptions",
			"moderator:manage:announcements",
			"moderator:manage:blocked_terms",
			"moderator:manage:chat_settings",
			"moderator:manage:shoutouts",
		},
		State:       state,
		ForceVerify: false,
	})

	fmt.Printf("Please authenticate at: %s\n", url)

	var tempauthcode string

	// Wait for a response from the OAuth2 provider
	select {
	case response := <-responseChan:
		if response.Error != "" {
			fmt.Printf("OAuth2 Error: %s", response.Error)
		} else {
			if authType == AuthToken {
				err := keys.AddKey("access-token", response.AccessToken)
				if err != nil {
					fmt.Printf("Failed to push access token to keystore: %s\n", err)
					return err
				}
				fmt.Printf("\nAccess token successfully received and pushed to keystore.\n")
				client.SetUserAccessToken(response.AccessToken)
			} else {
				tempauthcode = response.AccessToken
			}
		}
	case <-time.After(5 * time.Minute):
		fmt.Printf("Timeout waiting for OAuth2 callback response\n")
		// Signal the server to shut down
		shutdownChan <- true
	}

	if authType == AuthCode {
		clientSecret, err := keys.GetKey("client-secret")
		if err != nil {
			fmt.Printf("Failed to get Client Secret from keystore: %s\n", err)
			return err
		}

		c, err := helix.NewClient(&helix.Options{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURI:  "http://localhost:3024/redirect",
		})
		if err != nil {
			fmt.Printf("Unable to create client for token generation: %s\n", err)
			return err
		}

		resp, err := c.RequestUserAccessToken(tempauthcode)
		if err != nil {
			fmt.Printf("Unable to request access token: %s\n", err)
			return err
		}
		if resp.StatusCode >= 300 {
			fmt.Printf("Status code was bad: %v\n", resp)
			return fmt.Errorf("check status code information")
		}

		err = keys.AddKey("refresh-token", resp.Data.RefreshToken)
		if err != nil {
			fmt.Printf("Failed to push refresh token to keystore: %s\n", err)
			return err
		}

		err = keys.AddKey("access-token", resp.Data.AccessToken)
		if err != nil {
			fmt.Printf("Failed to push access token to keystore: %s\n", err)
			return err
		}

		fmt.Printf("\nAccess token and refresh token successfully received and pushed to keystore.\n")

		client.SetUserAccessToken(resp.Data.AccessToken)
	}

	return nil // This is unreachable most of the time?
}

// Authenticate is the process to get a user token from Twitch.
// It requires that the setup process (putting client ID into the keystore) has already taken place.
func RefreshToken(clientID string, clientSecret string, refreshToken string) (helix.Client, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})

	resp, err := client.RefreshUserAccessToken(refreshToken)
	if err != nil {
		fmt.Printf("Failed to refresh user access token: %s\n", err)
		return *client, err
	}

	err = keys.AddKey("refresh-token", resp.Data.RefreshToken)
	if err != nil {
		fmt.Printf("Failed to push refresh token to keystore: %s\n", err)
		return *client, err
	}

	err = keys.AddKey("access-token", resp.Data.AccessToken)
	if err != nil {
		fmt.Printf("Failed to push access token to keystore: %s\n", err)
		return *client, err
	}

	client.SetUserAccessToken(resp.Data.AccessToken)

	return *client, nil
}
