package cmd

import (
	"fmt"

	"github.com/monktype/msc/callback"
	"github.com/spf13/cobra"
)

// The version string is set by Git from the tag when building.
var version string

var (
	rootCmd = &cobra.Command{
		Use:   "msc",
		Short: "msc is a command line tool for interacting with the Twitch API",
		Long: `Monktype's Stream Commands is a command line tool for interacting
with the Twitch API.`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

// Like... I could separate these, but I can't be bothered right now.
func init() {
	// --- Before other flags, get the global flags read and set. ---
	var callbackPort int
	rootCmd.PersistentFlags().IntVar(&callbackPort, "callback-port", 3024, "Twitch->msc authentication callback port if default can't be used")

	// Set the PreRun to update the CallbackPort
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		callback.CallbackPort = callbackPort
	}

	// --- Call-specific flags now ---
	rootCmd.AddCommand(versionCmd)
	setupCmd.Flags().BoolP("no-auth", "n", false, "Skip trying to authenticate after running setup.")
	setupCmd.Flags().BoolP("secret", "s", false, "Add a secret for code authentication instead of token authentication.")
	setupCmd.Flags().StringP("client-id", "i", "", "Client ID from Twitch Dev portal")
	setupCmd.MarkFlagRequired("client-id")
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(userIDCmd)
	pollCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	pollCmd.MarkFlagRequired("channel-name")
	pollCmd.Flags().StringP("title", "t", "", "Title for poll")
	pollCmd.MarkFlagRequired("title")
	pollCmd.Flags().IntP("duration", "d", 0, "Duration in seconds")
	pollCmd.MarkFlagRequired("duration")
	pollCmd.Flags().BoolP("no-watch", "n", false, "Skip watching until the end of the poll, just print the poll ID.")
	pollCmd.Flags().BoolP("send-announcement", "a", false, "Send an announcement when the poll starts. 'New poll for X seconds: \"Poll Title\"'")
	pollCmd.Flags().BoolP("send-announcement-result", "A", false, "Send an announcement when the poll starts AND when the poll ends with the result. Implies -a.")
	rootCmd.AddCommand(pollCmd)
	announcementCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	announcementCmd.MarkFlagRequired("channel-name")
	announcementCmd.Flags().StringP("border-color", "b", "primary", "Border color (primary, blue, green, orange, purple)")
	rootCmd.AddCommand(announcementCmd)
	shoutoutCmd.Flags().StringP("channel-name", "c", "", "Channel name to send shoutout")
	shoutoutCmd.MarkFlagRequired("channel-name")
	shoutoutCmd.Flags().StringP("shoutout-name", "s", "", "Shoutout name")
	shoutoutCmd.MarkFlagRequired("shoutout-name")
	rootCmd.AddCommand(shoutoutCmd)
	startadCmd.Flags().StringP("channel-name", "c", "", "Channel name to start ads")
	startadCmd.MarkFlagRequired("channel-name")
	startadCmd.Flags().IntP("length", "l", 60, "Ad length in seconds (30, 60, 90, 120, 150, 180)")
	startadCmd.MarkFlagRequired("length")
	rootCmd.AddCommand(startadCmd)
	rootCmd.AddCommand(emoteonlyCmd)
	emoteonlyOnCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	emoteonlyOnCmd.MarkFlagRequired("channel-name")
	emoteonlyCmd.AddCommand(emoteonlyOnCmd)
	emoteonlyOffCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	emoteonlyOffCmd.MarkFlagRequired("channel-name")
	emoteonlyCmd.AddCommand(emoteonlyOffCmd)
	rootCmd.AddCommand(followeronlyCmd)
	followeronlyOnCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	followeronlyOnCmd.MarkFlagRequired("channel-name")
	followeronlyCmd.AddCommand(followeronlyOnCmd)
	followeronlyOffCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	followeronlyOffCmd.MarkFlagRequired("channel-name")
	followeronlyCmd.AddCommand(followeronlyOffCmd)
	followeronlyDurationCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	followeronlyDurationCmd.MarkFlagRequired("channel-name")
	followeronlyDurationCmd.Flags().IntP("duration", "d", 0, "Duration in minutes (0..129600)")
	followeronlyDurationCmd.MarkFlagRequired("duration")
	followeronlyCmd.AddCommand(followeronlyDurationCmd)
	rootCmd.AddCommand(slowmodeCmd)
	slowmodeOnCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	slowmodeOnCmd.MarkFlagRequired("channel-name")
	slowmodeCmd.AddCommand(slowmodeOnCmd)
	slowmodeOffCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	slowmodeOffCmd.MarkFlagRequired("channel-name")
	slowmodeCmd.AddCommand(slowmodeOffCmd)
	slowmodeDurationCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	slowmodeDurationCmd.MarkFlagRequired("channel-name")
	slowmodeDurationCmd.Flags().IntP("duration", "d", 0, "Duration in seconds (3..120)")
	slowmodeDurationCmd.MarkFlagRequired("duration")
	slowmodeCmd.AddCommand(slowmodeDurationCmd)
	rootCmd.AddCommand(submodeCmd)
	submodeOnCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	submodeOnCmd.MarkFlagRequired("channel-name")
	submodeCmd.AddCommand(submodeOnCmd)
	submodeOffCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	submodeOffCmd.MarkFlagRequired("channel-name")
	submodeCmd.AddCommand(submodeOffCmd)
	rootCmd.AddCommand(rewardsCmd)
	rewardscreateCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	rewardscreateCmd.MarkFlagRequired("channel-name")
	rewardscreateCmd.Flags().StringP("title", "t", "", "Reward title")
	rewardscreateCmd.MarkFlagRequired("title")
	rewardscreateCmd.Flags().IntP("points", "p", 0, "The point cost of the item")
	rewardscreateCmd.MarkFlagRequired("points")
	rewardscreateCmd.Flags().StringP("background-color", "b", "", "Background color in hex notation (ie #9147FF)")
	rewardscreateCmd.Flags().StringP("user-prompt", "u", "", "Set prompt (shows if input is enabled)")
	rewardscreateCmd.Flags().BoolP("input-required", "i", false, "Make user input required")
	rewardsCmd.AddCommand(rewardscreateCmd)
	rewardsdeleteCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	rewardsdeleteCmd.MarkFlagRequired("channel-name")
	rewardsdeleteCmd.Flags().StringP("reward", "r", "", "Reward ID (UUID)")
	rewardsdeleteCmd.MarkFlagRequired("reward")
	rewardsCmd.AddCommand(rewardsdeleteCmd)
	rewardsgetCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	rewardsgetCmd.MarkFlagRequired("channel-name")
	rewardsCmd.AddCommand(rewardsgetCmd)
	rewardsredemptionsCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	rewardsredemptionsCmd.MarkFlagRequired("channel-name")
	rewardsredemptionsCmd.Flags().StringP("reward", "r", "", "Reward ID (UUID)")
	rewardsredemptionsCmd.MarkFlagRequired("reward")
	rewardsredemptionsCmd.Flags().StringP("status", "s", "", "Redemption status. Can only be CANCELED, FULFILLED, and UNFULFILLED")
	rewardsredemptionsCmd.MarkFlagRequired("status")
	rewardsCmd.AddCommand(rewardsredemptionsCmd)
	rewardscancelCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	rewardscancelCmd.MarkFlagRequired("channel-name")
	rewardscancelCmd.Flags().StringP("reward", "r", "", "Reward ID (UUID). This is the button on the Twitch UI, not the redeemed reward instance.")
	rewardscancelCmd.MarkFlagRequired("reward")
	rewardscancelCmd.Flags().StringP("redemption", "i", "", "Redemption ID (UUID). This is the redeemed reward instance.")
	rewardscancelCmd.MarkFlagRequired("redemption")
	rewardsCmd.AddCommand(rewardscancelCmd)
	rewardsfulfillCmd.Flags().StringP("channel-name", "c", "", "Target channel name")
	rewardsfulfillCmd.MarkFlagRequired("channel-name")
	rewardsfulfillCmd.Flags().StringP("reward", "r", "", "Reward ID (UUID). This is the button on the Twitch UI, not the redeemed reward instance.")
	rewardsfulfillCmd.MarkFlagRequired("reward")
	rewardsfulfillCmd.Flags().StringP("redemption", "i", "", "Redemption ID (UUID). This is the redeemed reward instance.")
	rewardsfulfillCmd.MarkFlagRequired("redemption")
	rewardsCmd.AddCommand(rewardsfulfillCmd)
	startApiCmd.Flags().IntP("port", "p", 8080, "Port for API server")
	rootCmd.AddCommand(startApiCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Monktype's Stream Commands",
	Run: func(cmd *cobra.Command, args []string) {
		// Print the version number
		if version == "" {
			fmt.Printf("It looks like this is a development build; no version tagged.\n")
		} else {
			// This can be set with `go build -ldflags "-X 'github.com/monktype/msc/cmd.version=${VERSION}'"` during build.
			fmt.Printf("Version: %s\n", version)
		}
	},
}
