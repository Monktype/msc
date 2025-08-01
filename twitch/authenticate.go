package twitch

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/monktype/msc/keys"
	"github.com/monktype/msc/server"
	"github.com/nicklaw5/helix/v2"
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
func Authenticate() error {
	// Generate a random state
	state, err := generateRandomState(16)
	if err != nil {
		fmt.Printf("Failed to generate random state: %s\n", err)
		return err
	}

	responseChan := make(chan server.CallbackResponse)
	shutdownChan := make(chan bool)

	go func() {
		err := server.StartCallbackServer(state, responseChan, shutdownChan)
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

	url := client.GetAuthorizationURL(&helix.AuthorizationURLParams{
		ResponseType: "token", // "token" for the implicit grant flow method.
		Scopes: []string{
			"channel:edit:commercial",
			"channel:manage:polls",
			"channel:manage:predictions",
			"channel:manage:redemptions",
			"moderator:manage:announcements",
		},
		State:       state,
		ForceVerify: false,
	})

	fmt.Printf("Please authenticate at: %s\n", url)

	// Wait for a response from the OAuth2 provider
	select {
	case response := <-responseChan:
		if response.Error != "" {
			fmt.Printf("OAuth2 Error: %s", response.Error)
		} else {
			err := keys.AddKey("access-token", response.AccessToken)
			if err != nil {
				fmt.Printf("Failed to push access token to keystore: %s\n", err)
				return err
			}
			fmt.Printf("Access token successfully received and pushed to keystore.\n")
		}
		// Signal the server to shut down
		shutdownChan <- true
	case <-time.After(5 * time.Minute):
		fmt.Printf("Timeout waiting for OAuth2 callback response\n")
		// Signal the server to shut down
		shutdownChan <- true
	}

	return nil // This is unreachable most of the time?
}
