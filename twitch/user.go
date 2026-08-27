package twitch

import (
	"fmt"
	"sync"

	"github.com/monktype/msc/keys"
	"github.com/nicklaw5/helix/v2"
)

var getClientLock sync.Mutex

// Create a Helix (Twitch) client, return the usable client struct (helix.Client) and error.
func GetClient() (*helix.Client, error) {
	secretPresent := false

	// This resolves a potential key refresh race condition when multiple clients reach the tool's API server
	// when a key refresh is necessary. It does, however, mean that only one GetClient instance can run at a time
	// which might add some wait time.
	// It does not resolve two or more instances of this tool running and finding an invalid key at the same time,
	// that's still a possible race condition.
	getClientLock.Lock()
	defer getClientLock.Unlock()

	// This is to check for the existence of "client-secret" in the keychain to decide what to do if the validation fails.
	clientsecret, err := keys.GetKey("client-secret")
	if err == nil && clientsecret != "" {
		secretPresent = true
	}

	clientID, err := keys.GetKey("client-id")
	if err != nil {
		fmt.Printf("Failed to get Client ID from keystore: %s\n", err)
		return nil, err
	}

	accessToken, err := keys.GetKey("access-token")
	if err != nil {
		fmt.Printf("Failed to get access token from keystore: %s\n", err)
		return nil, err
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:        clientID,
		UserAccessToken: accessToken,
	})
	if err != nil {
		fmt.Printf("Failed to create Helix client: %s\n", err)
		return nil, err
	}

	isValid, resp, err := client.ValidateToken(accessToken)
	if err != nil {
		fmt.Printf("Token validation failed: %s\n", err)
		return nil, err
	}

	if isValid == false || (isValid == true && resp.Data.ExpiresIn < 330 && secretPresent) {
		if secretPresent {
			refreshToken, err := keys.GetKey("refresh-token")
			if err != nil {
				fmt.Printf("Failed to get refresh token from keystore: %s\n", err)
				if !isValid {
					fmt.Printf("Could not use refresh token to update code flow.\nTry `msc authenticate` again.\n")
					return nil, err
				}
				fmt.Printf("Continuing for now, but your token expires soon.\n")
				return client, nil
			}
			refreshedclient, err := RefreshToken(clientID, clientsecret, refreshToken)
			if err != nil {
				fmt.Printf("Failed to refresh auth token: %s\n", err)
				if isValid {
					fmt.Printf("Continuing for now after refresh failure; your token expires soon.\n")
					return client, nil
				}
				return nil, err
			}
			return refreshedclient, nil
		} else { // Presumed token access that's expired.
			fmt.Printf("Token expired. Run `msc authenticate` to re-authenticate.\n")
			return nil, fmt.Errorf("token expired, re-authenticate")
		}

	}

	return client, nil
}

// GetUserID gets User ID from a username.
// Takes helix.Client and username string.
// Returns ID as string, error.
func GetUserID(c *helix.Client, username string) (string, error) {
	resp, err := c.GetUsers(&helix.UsersParams{
		Logins: []string{username},
	})
	if err != nil {
		fmt.Printf("Failed to get user %s: %s\n", username, err)
		return "", err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return "", fmt.Errorf("check status code information")
	}
	if len(resp.Data.Users) == 0 {
		fmt.Printf("User %s was not found.\n", username)
		return "", fmt.Errorf("user %s not found", username)
	}

	return resp.Data.Users[0].ID, nil
}

// GetMyUserID gets the User ID from the current user.
// Takes helix.Client.
// Returns ID as string, error.
func GetMyUserID(c *helix.Client) (string, error) {
	resp, err := c.GetUsers(&helix.UsersParams{}) // the magic is not sending any parameters
	if err != nil {
		fmt.Printf("Failed to get my current user: %s\n", err)
		return "", err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return "", fmt.Errorf("check status code information")
	}
	if len(resp.Data.Users) == 0 {
		fmt.Printf("Failed to resolve the current user's ID: response was empty.\n")
		return "", fmt.Errorf("could not resolve the current user's ID")
	}

	return resp.Data.Users[0].ID, nil
}
