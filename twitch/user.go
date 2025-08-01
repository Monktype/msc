package twitch

import (
	"fmt"

	"github.com/monktype/msc/keys"
	"github.com/nicklaw5/helix/v2"
)

// Create a Helix (Twitch) client, return the usable client struct (helix.Client) and error.
func GetClient() (helix.Client, error) {
	clientID, err := keys.GetKey("client-id")
	if err != nil {
		fmt.Printf("Failed to get Client ID from keystore: %s\n", err)
		return helix.Client{}, err
	}

	accessToken, err := keys.GetKey("access-token")
	if err != nil {
		fmt.Printf("Failed to get access token from keystore: %s\n", err)
		return helix.Client{}, err
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:        clientID,
		UserAccessToken: accessToken,
	})
	if err != nil {
		fmt.Printf("Failed to create Helix client: %s\n", err)
		return helix.Client{}, err
	}

	return *client, nil
}

// GetUserID gets User ID from a username.
// Takes helix.Client and username string.
// Returns ID as string, error.
func GetUserID(c helix.Client, username string) (string, error) {
	resp, err := c.GetUsers(&helix.UsersParams{
		Logins: []string{username},
	})
	if err != nil {
		fmt.Printf("Failed to get user %s: %s\n", username, err)
		return "", err
	}

	return resp.Data.Users[0].ID, nil
}

// GetMyUserID gets the User ID from the current user.
// Takes helix.Client.
// Returns ID as string, error.
func GetMyUserID(c helix.Client) (string, error) {
	resp, err := c.GetUsers(&helix.UsersParams{}) // the magic is not sending any parameters
	if err != nil {
		fmt.Printf("Failed to get my current user: %s\n", err)
		return "", err
	}

	return resp.Data.Users[0].ID, nil
}
