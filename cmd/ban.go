package cmd

import (
	"fmt"
	"time"

	"github.com/monktype/msc/twitch"
	"github.com/nicklaw5/helix/v2"
	"github.com/spf13/cobra"
)

var banCmd = &cobra.Command{
	Use:   "ban",
	Short: "Ban and manage banned/timed-out users",
}

var banListCmd = &cobra.Command{
	Use:   "list",
	Short: "List banned and timed-out users",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		after, err := cmd.Flags().GetString("after")
		if err != nil {
			return err
		}

		channelid, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		listAll, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}

		var bans []helix.Ban
		if listAll {
			bans, err = twitch.GetAllBannedUsers(c, channelid)
			if err != nil {
				return err
			}
		} else {
			page, err := twitch.GetBannedUsers(c, channelid, after)
			if err != nil {
				return err
			}
			bans = page.Bans
		}

		if len(bans) != 0 {
			fmt.Printf("Current banned/timed-out users on channel:\n\n")

			for _, ban := range bans {
				fmt.Printf("%s (%s) [%s]\n", ban.UserLogin, ban.UserName, ban.UserID)
				if ban.ExpiresAt.IsZero() {
					fmt.Printf("  perma-ban\n")
				} else {
					fmt.Printf("  timeout until %s\n", ban.ExpiresAt.Format(time.RFC3339))
				}
			}

			fmt.Printf("\n")
		} else {
			fmt.Printf("No banned/timed-out users found.\n")
		}

		return nil
	},
}

var banAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Ban or time out a user",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		username, err := cmd.Flags().GetString("user")
		if err != nil {
			return err
		}

		userid, err := cmd.Flags().GetString("user-id")
		if err != nil {
			return err
		}

		duration, err := cmd.Flags().GetInt("duration")
		if err != nil {
			return err
		}

		reason, err := cmd.Flags().GetString("reason")
		if err != nil {
			return err
		}

		channelid, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		var targetid string
		switch {
		case userid != "":
			targetid = userid
		case username != "":
			targetid, err = twitch.GetUserID(c, username)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("provide -u <username> or --user-id <id>")
		}

		moderatorid, err := twitch.GetMyUserID(c)
		if err != nil {
			return err
		}

		err = twitch.BanUser(c, moderatorid, channelid, targetid, duration, reason)
		if err != nil {
			return err
		}

		return nil
	},
}

var banRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a ban or timeout",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		username, err := cmd.Flags().GetString("user")
		if err != nil {
			return err
		}

		userid, err := cmd.Flags().GetString("user-id")
		if err != nil {
			return err
		}

		channelid, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		var targetid string
		switch {
		case userid != "":
			targetid = userid
		case username != "":
			targetid, err = twitch.GetUserID(c, username)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("provide -u <username> or --user-id <id>")
		}

		moderatorid, err := twitch.GetMyUserID(c)
		if err != nil {
			return err
		}

		err = twitch.UnbanUser(c, moderatorid, channelid, targetid)
		if err != nil {
			return err
		}

		return nil
	},
}
