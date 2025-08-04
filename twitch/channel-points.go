package twitch

import (
	"fmt"

	"github.com/nicklaw5/helix/v2"
)

// This is here for information because it's so big, but use the helix version.
/*
type ChannelCustomRewardsParams struct {
    BroadcasterID                     string
	Title                             string
	Cost                              int // minimum is 1
	Prompt                            string // If IsUserInputRequired = True; 200 chat max.
	IsEnabled                         bool
	BackgroundColor                   string // Hex string (example: #9147FF)
	IsUserInputRequired               bool
	IsMaxPerStreamEnabled             bool
	MaxPerStream                      int // Minimum if IsMaxPerStreamEnabled = True is 1.
	IsMaxPerUserPerStreamEnabled      bool
	MaxPerUserPerStream               int // Minimum if IsMaxPerUserPerStreamEnabled = True is 1.
	IsGlobalCooldownEnabled           bool
	GlobalCooldownSeconds             int // Minimum value 1, but minimum value 60 for it to be shown on Twitch UX.
	ShouldRedemptionsSkipRequestQueue bool
}
*/

// CreateReward creates a Twitch custom channel points reward with the given ChannelCustomRewardsParams.
// Returns error.
func CreateReward(c helix.Client, params helix.ChannelCustomRewardsParams) (string, error) {
	resp, err := c.CreateCustomReward(&params)
	if err != nil {
		fmt.Printf("Creating a reward failed: %s\n", err)
		return "", err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return "", fmt.Errorf("check status code information")
	}

	fmt.Printf("Created custom reward with ID %s\n", resp.Data.ChannelCustomRewards[0].ID)
	return resp.Data.ChannelCustomRewards[0].ID, nil
}

// DeleteReward deletes a Twitch custom channel points reward with the given channel ID and reward ID.
// Returns error.
func DeleteReward(c helix.Client, channelID string, rewardID string) error {
	resp, err := c.DeleteCustomRewards(&helix.DeleteCustomRewardsParams{
		BroadcasterID: channelID,
		ID:            rewardID,
	})
	if err != nil {
		fmt.Printf("Deleting a reward failed: %s\n", err)
		return err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return fmt.Errorf("check status code information")
	}

	fmt.Printf("Deleted reward with ID %s\n", rewardID)
	return nil
}

// GetRewards gets Twitch custom channel points rewards for the given channel ID.
// Returns []helix.ChannelCustomReward and error.
func GetRewards(c helix.Client, channelID string) ([]helix.ChannelCustomReward, error) {
	var emptyRewards []helix.ChannelCustomReward

	resp, err := c.GetCustomRewards(&helix.GetCustomRewardsParams{
		BroadcasterID: channelID,
	})
	if err != nil {
		fmt.Printf("Getting channel rewards failed: %s\n", err)
		return emptyRewards, err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return emptyRewards, fmt.Errorf("check status code information")
	}

	return resp.Data.ChannelCustomRewards, nil
}

// GetRedemptions gets Twitch custom channel points rewards' redemptions for the given channel ID, reward ID, and status.
// Returns []helix.ChannelCustomRewardsRedemption and error.
func GetRedemptions(c helix.Client, channelID string, rewardID string, status string) ([]helix.ChannelCustomRewardsRedemption, error) {
	var emptyRedemption []helix.ChannelCustomRewardsRedemption

	resp, err := c.GetCustomRewardsRedemptions(&helix.GetCustomRewardsRedemptionsParams{
		BroadcasterID: channelID,
		RewardID:      rewardID,
		Status:        status,
	})
	if err != nil {
		fmt.Printf("Getting channel reward redeptions failed: %s\n", err)
		return emptyRedemption, err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return emptyRedemption, fmt.Errorf("check status code information")
	}

	return resp.Data.Redemptions, nil
}

// CancelRedemption cancels a Twitch custom channel points rewards' redemption for the given channel ID, reward ID, and redemption ID.
// Returns []helix.ChannelCustomRewardsRedemption and error.
func CancelRedemption(c helix.Client, channelID string, rewardID string, redemptionID string) ([]helix.ChannelCustomRewardsRedemption, error) {
	var emptyRedemption []helix.ChannelCustomRewardsRedemption

	resp, err := c.UpdateChannelCustomRewardsRedemptionStatus(&helix.UpdateChannelCustomRewardsRedemptionStatusParams{
		BroadcasterID: channelID,
		RewardID:      rewardID,
		ID:            redemptionID,
		Status:        "CANCELED",
	})
	if err != nil {
		fmt.Printf("Cancelling a channel reward redeption failed: %s\n", err)
		return emptyRedemption, err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return emptyRedemption, fmt.Errorf("check status code information")
	}

	return resp.Data.Redemptions, nil
}

// FulfillRedemption fulfills a Twitch custom channel points rewards' redemption for the given channel ID, reward ID, and redemption ID.
// Returns []helix.ChannelCustomRewardsRedemption and error.
func FulfillRedemption(c helix.Client, channelID string, rewardID string, redemptionID string) ([]helix.ChannelCustomRewardsRedemption, error) {
	var emptyRedemption []helix.ChannelCustomRewardsRedemption

	resp, err := c.UpdateChannelCustomRewardsRedemptionStatus(&helix.UpdateChannelCustomRewardsRedemptionStatusParams{
		BroadcasterID: channelID,
		RewardID:      rewardID,
		ID:            redemptionID,
		Status:        "FULFILLED",
	})
	if err != nil {
		fmt.Printf("Fulfilling a channel reward redeption failed: %s\n", err)
		return emptyRedemption, err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return emptyRedemption, fmt.Errorf("check status code information")
	}

	return resp.Data.Redemptions, nil
}
