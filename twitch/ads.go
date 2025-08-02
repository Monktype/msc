package twitch

import (
	"fmt"

	"github.com/nicklaw5/helix/v2"
)

func StartCommercial(c helix.Client, channelID string, length helix.AdLengthEnum) error {
	resp, err := c.StartCommercial(&helix.StartCommercialParams{
		BroadcasterID: channelID,
		Length:        length,
	})
	if err != nil {
		fmt.Printf("Commercial failed to start: %s\n", err)
		return err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return fmt.Errorf("check status code information")
	}

	return nil
}
