package cmd

import (
	"fmt"

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

func init() {
	rootCmd.AddCommand(versionCmd)
	setupCmd.Flags().BoolP("no-auth", "n", false, "Skip trying to authenticate after running setup.")
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
