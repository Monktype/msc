package cmd

import (
	"fmt"

	"github.com/monktype/msc/twitch"
	"github.com/spf13/cobra"
)

var userIDCmd = &cobra.Command{
	Use:   "userid",
	Short: "Look up user ID from username",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		userID, err := twitch.GetUserID(c, args[0])
		if err != nil {
			return err
		}

		fmt.Printf("Username %s = ID %s\n", args[0], userID)
		return nil
	},
}
