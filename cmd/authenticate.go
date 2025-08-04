package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

		addsecret, err := cmd.Flags().GetBool("secret")
		if err != nil {
			return err
		}

		if addsecret {
			fmt.Printf("Please input your client secret here -> ")
			reader := bufio.NewReader(os.Stdin)
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("Failed to read your input: %s\n", err)
				return fmt.Errorf("failed to read client-secret from stdin")
			}
			cleanedLine := strings.TrimSpace(line)

			err = keys.AddKey("client-secret", cleanedLine)
			if err != nil {
				return err
			}

			fmt.Printf("\n")
		}

		err = keys.AddKey("client-id", cid)
		if err != nil {
			return err
		}

		if !skipauth {
			authtype := twitch.AuthToken
			if addsecret {
				authtype = twitch.AuthCode
			}
			err = twitch.Authenticate(authtype)
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
		authtype := twitch.AuthToken

		// This is to check for the existence of "client-secret" in the keychain to decide what type of authentication to use.
		clientsecret, err := keys.GetKey("client-secret")
		if err == nil && clientsecret != "" {
			authtype = twitch.AuthCode
		}

		err = twitch.Authenticate(authtype)
		if err != nil {
			return err
		}

		return nil
	},
}
