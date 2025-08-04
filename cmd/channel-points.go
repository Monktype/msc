package cmd

import (
	"fmt"
	"strings"

	"github.com/monktype/msc/twitch"
	"github.com/nicklaw5/helix/v2"
	"github.com/spf13/cobra"
)

var rewardsCmd = &cobra.Command{
	Use:   "reward",
	Short: "Custom Channel Point Reward",
}

var rewardscreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Channel Point Reward",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		title, err := cmd.Flags().GetString("title")
		if err != nil {
			return err
		}

		cost, err := cmd.Flags().GetInt("points")
		if err != nil {
			return err
		}

		colorcode, err := cmd.Flags().GetString("background-color")
		if err != nil {
			return err
		}

		userinputrequired, err := cmd.Flags().GetBool("input-required")
		if err != nil {
			return err
		}

		userprompt, err := cmd.Flags().GetString("user-prompt")
		if err != nil {
			return err
		}

		channelid, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		var params helix.ChannelCustomRewardsParams

		params.BroadcasterID = channelid
		params.Title = title
		params.Cost = cost
		params.IsEnabled = true // That's the point of this function
		params.BackgroundColor = colorcode
		params.IsUserInputRequired = userinputrequired
		params.Prompt = userprompt // I'm leaving it possible to set this without enabling user input.

		_, err = twitch.CreateReward(c, params)
		if err != nil {
			return err
		}

		return nil
	},
}

var rewardsdeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Channel Point Reward",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		rewardid, err := cmd.Flags().GetString("reward")
		if err != nil {
			return err
		}

		channelid, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		err = twitch.DeleteReward(c, channelid, rewardid)
		if err != nil {
			return err
		}

		return nil
	},
}

var rewardsgetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get Channel Point Rewards",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		channelid, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		rewards, err := twitch.GetRewards(c, channelid)
		if err != nil {
			return err
		}

		if len(rewards) != 0 {
			fmt.Printf("Current rewards on channel:\n\n")

			for _, reward := range rewards {
				fmt.Printf("%s:\t%s (%d)\n", reward.ID, reward.Title, reward.Cost)
			}

			fmt.Printf("\n")
		} else {
			fmt.Printf("No rewards currently found.\n")
		}

		return nil
	},
}

var rewardsredemptionsCmd = &cobra.Command{
	Use:   "redemptions",
	Short: "Get Channel Point Reward Redemptions",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		rewardid, err := cmd.Flags().GetString("reward")
		if err != nil {
			return err
		}

		statusraw, err := cmd.Flags().GetString("status")
		if err != nil {
			return err
		}

		status := strings.ToUpper(statusraw) // The Twitch API wants it upper-cased.
		if status != "CANCELED" && status != "FULFILLED" && status != "UNFULFILLED" {
			fmt.Printf("Status can only be CANCELED or FULFILLED or UNFULFILLED.\n")
			return fmt.Errorf("status string is incorrect")
		}

		channelid, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		redemptions, err := twitch.GetRedemptions(c, channelid, rewardid, status)
		if err != nil {
			return err
		}

		if len(redemptions) != 0 {
			fmt.Printf("Current redemptions on channel:\n\n")

			for _, redemption := range redemptions {
				fmt.Printf("%s by user %s (%s) (%s)\n", redemption.ID, redemption.UserName, redemption.UserID, redemption.Status)
			}

			fmt.Printf("\n")
		} else {
			fmt.Printf("No rewards currently found.\n")
		}

		return nil
	},
}

var rewardscancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel a Channel Point Reward Redemption",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		rewardid, err := cmd.Flags().GetString("reward")
		if err != nil {
			return err
		}

		redemptionid, err := cmd.Flags().GetString("redemption")
		if err != nil {
			return err
		}

		channelid, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		redemptions, err := twitch.CancelRedemption(c, channelid, rewardid, redemptionid)
		if err != nil {
			return err
		}

		if len(redemptions) != 0 {
			fmt.Printf("Returned redemption after operation:\n\n")

			for _, redemption := range redemptions {
				fmt.Printf("%s by user %s (%s) (%s)\n", redemption.ID, redemption.UserName, redemption.UserID, redemption.Status)
			}

			fmt.Printf("\n")
		}

		return nil
	},
}

var rewardsfulfillCmd = &cobra.Command{
	Use:   "fulfill",
	Short: "Fulfill a Channel Point Reward Redemption",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		rewardid, err := cmd.Flags().GetString("reward")
		if err != nil {
			return err
		}

		redemptionid, err := cmd.Flags().GetString("redemption")
		if err != nil {
			return err
		}

		channelid, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		redemptions, err := twitch.FulfillRedemption(c, channelid, rewardid, redemptionid)
		if err != nil {
			return err
		}

		if len(redemptions) != 0 {
			fmt.Printf("Returned redemption after operation:\n\n")

			for _, redemption := range redemptions {
				fmt.Printf("%s by user %s (%s) (%s)\n", redemption.ID, redemption.UserName, redemption.UserID, redemption.Status)
			}

			fmt.Printf("\n")
		}

		return nil
	},
}
