package cmd

import (
	"fmt"

	"github.com/monktype/msc/twitch"
	"github.com/spf13/cobra"
)

var pollCmd = &cobra.Command{
	Use:   "poll",
	Short: "Create a poll with -c (channel name), -d (duration in seconds), -t (title), followed by options",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			fmt.Printf("At least 2 poll options are required, only %d provided.\n", len(args))
			return fmt.Errorf("give the right number of arguments")
		}
		if len(args) > 10 {
			fmt.Printf("At most 10 poll options are allowed, %d provided.\n", len(args))
			return fmt.Errorf("give the right number of arguments")
		}

		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		title, err := cmd.Flags().GetString("title")
		if err != nil {
			return err
		}

		duration, err := cmd.Flags().GetInt("duration")
		if err != nil {
			return err
		}

		nowatch, err := cmd.Flags().GetBool("no-watch")
		if err != nil {
			return err
		}

		sendannouncement, err := cmd.Flags().GetBool("send-announcement")
		if err != nil {
			return err
		}

		sendannouncementresult, err := cmd.Flags().GetBool("send-announcement-result")
		if err != nil {
			return err
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		userID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		pollID, err := twitch.CreatePoll(c, userID, title, duration, args)
		if err != nil {
			return err
		}

		if sendannouncement || sendannouncementresult {
			myUserID, err := twitch.GetMyUserID(c)
			if err != nil {
				// This isn't a fatal thing, just mention it and skip the announcement part.
				fmt.Printf("Failed to send announcement because getting the current user failed, but continuing with the poll: %s\n", err)
			} else {
				err = twitch.SendAnnouncement(c, myUserID, userID, twitch.AnnouncementColorPrimary, fmt.Sprintf("New poll for %d seconds! \"%s\"", duration, title))
				if err != nil {
					// This isn't a fatal thing, just mention it.
					fmt.Printf("Failed to send announcement but continuing with the poll: %s\n", err)
				}
			}
		}

		if nowatch {
			return nil
		}

		fmt.Printf("Waiting for poll completion...\n")
		resultstring, err := twitch.WatchPollCompletion(c, userID, pollID)
		if err != nil {
			return err
		}

		fmt.Printf("\n%s\n", resultstring)

		if sendannouncementresult {
			myUserID, err := twitch.GetMyUserID(c)
			if err != nil {
				// This isn't a fatal thing, just mention it and skip the announcement part.
				fmt.Printf("Failed to send announcement result because getting the current user failed, but continuing with the poll: %s\n", err)
			} else {
				err = twitch.SendAnnouncement(c, myUserID, userID, twitch.AnnouncementColorPrimary, fmt.Sprintf("Poll \"%s\" finished: %s", title, resultstring))
				if err != nil {
					// This isn't a fatal thing, just mention it.
					fmt.Printf("Failed to send announcement result but continuing with the poll: %s\n", err)
				}
			}
		}

		return nil
	},
}
