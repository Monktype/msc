package cmd

import (
	"github.com/monktype/msc/keys"
	"github.com/monktype/msc/twitch"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Run the setup process",
	RunE: func(cmd *cobra.Command, args []string) error {
		cid, err := cmd.Flags().GetString("client-id")
		if err != nil {
			return err
		}

		skipauth, err := cmd.Flags().GetBool("no-auth")
		if err != nil {
			return err
		}

		err = keys.AddKey("client-id", cid)
		if err != nil {
			return err
		}

		if !skipauth {
			err = twitch.Authenticate()
			if err != nil {
				return err
			}
		}

		return nil
	},
}

var authCmd = &cobra.Command{
	Use:   "authenticate",
	Short: "Run the authentication process",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := twitch.Authenticate()
		if err != nil {
			return err
		}

		return nil
	},
}
