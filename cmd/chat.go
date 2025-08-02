package cmd

import (
	"fmt"
	"strings"

	"github.com/monktype/msc/twitch"
	"github.com/spf13/cobra"
)

var announcementCmd = &cobra.Command{
	Use:   "announcement",
	Short: "Create an announcement with -c (channel name), -b (border-color), followed by announcement message",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			fmt.Printf("At least 1 word is required.\n")
			return fmt.Errorf("")
		}

		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		userColor, err := cmd.Flags().GetString("border-color")
		if err != nil {
			return err
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		userID, err := twitch.GetMyUserID(c)
		if err != nil {
			return err
		}

		var selectedColor twitch.AnnouncementColor
		foundColor := false

		for color, name := range twitch.AnnouncementColorMap {
			if strings.EqualFold(userColor, name) {
				selectedColor = color
				foundColor = true
				break
			}
		}

		if !foundColor {
			fmt.Printf("Color %s is invalid; please select \"primary\", \"blue\", \"green\", \"orange\", or \"purple\"", userColor)
			return fmt.Errorf("please select a valid color or let it default to \"primary\"")
		}

		err = twitch.SendAnnouncement(c, userID, channelID, selectedColor, strings.Join(args, " "))
		if err != nil {
			return err
		}

		return nil
	},
}

var shoutoutCmd = &cobra.Command{
	Use:   "shoutout",
	Short: "Create a shoutout with -c (channel name), -s (shoutout name)",
	RunE: func(cmd *cobra.Command, args []string) error {
		shoutoutname, err := cmd.Flags().GetString("shoutout-name")
		if err != nil {
			return err
		}

		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		targetID, err := twitch.GetUserID(c, shoutoutname)
		if err != nil {
			return err
		}

		userID, err := twitch.GetMyUserID(c)
		if err != nil {
			return err
		}

		err = twitch.SendShoutout(c, userID, channelID, targetID)
		if err != nil {
			return err
		}

		return nil
	},
}

var emoteonlyCmd = &cobra.Command{
	Use:   "emote-only",
	Short: "Enable or Disable emote-only mode with -c (channel name) flag",
}

var emoteonlyOnCmd = &cobra.Command{
	Use:   "on",
	Short: "Enable emote-only mode with -c (channel name) flag",
	RunE: func(cmd *cobra.Command, args []string) error {
		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		userID, err := twitch.GetMyUserID(c)
		if err != nil {
			return err
		}

		err = twitch.EmoteOnly(c, userID, channelID, true)
		if err != nil {
			return err
		}

		return nil
	},
}

var emoteonlyOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Disable emote-only mode with -c (channel name) flag",
	RunE: func(cmd *cobra.Command, args []string) error {
		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		userID, err := twitch.GetMyUserID(c)
		if err != nil {
			return err
		}

		err = twitch.EmoteOnly(c, userID, channelID, false)
		if err != nil {
			return err
		}

		return nil
	},
}

var followeronlyCmd = &cobra.Command{
	Use:   "follower-only",
	Short: "Enable or Disable follower-only mode with -c (channel name) flag",
}

var followeronlyOnCmd = &cobra.Command{
	Use:   "on",
	Short: "Enable follower-only mode with -c (channel name) flag",
	RunE: func(cmd *cobra.Command, args []string) error {
		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		userID, err := twitch.GetMyUserID(c)
		if err != nil {
			return err
		}

		err = twitch.FollowerOnly(c, userID, channelID, true)
		if err != nil {
			return err
		}

		return nil
	},
}

var followeronlyOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Disable follower-only mode with -c (channel name) flag",
	RunE: func(cmd *cobra.Command, args []string) error {
		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		userID, err := twitch.GetMyUserID(c)
		if err != nil {
			return err
		}

		err = twitch.FollowerOnly(c, userID, channelID, false)
		if err != nil {
			return err
		}

		return nil
	},
}

var followeronlyDurationCmd = &cobra.Command{
	Use:   "duration",
	Short: "Set follower-only mode with -c (channel name) and -d (duration in minutes) flags. Duration can be 0..129600 minutes.",
	RunE: func(cmd *cobra.Command, args []string) error {
		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		duration, err := cmd.Flags().GetInt("duration")
		if err != nil {
			return err
		}

		if duration < 0 || duration > 129600 {
			fmt.Printf("Duration in minutes can only be between 0 and 129600.\n")
			return fmt.Errorf("use a valid duration")
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		userID, err := twitch.GetMyUserID(c)
		if err != nil {
			return err
		}

		err = twitch.FollowerOnlyDuration(c, userID, channelID, duration)
		if err != nil {
			return err
		}

		return nil
	},
}

var slowmodeCmd = &cobra.Command{
	Use:   "slowmode",
	Short: "Enable or Disable slowmode with -c (channel name) flag",
}

var slowmodeOnCmd = &cobra.Command{
	Use:   "on",
	Short: "Enable slowmode mode with -c (channel name) flag",
	RunE: func(cmd *cobra.Command, args []string) error {
		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		userID, err := twitch.GetMyUserID(c)
		if err != nil {
			return err
		}

		err = twitch.Slowmode(c, userID, channelID, true)
		if err != nil {
			return err
		}

		return nil
	},
}

var slowmodeOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Disable slowmode with -c (channel name) flag",
	RunE: func(cmd *cobra.Command, args []string) error {
		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		userID, err := twitch.GetMyUserID(c)
		if err != nil {
			return err
		}

		err = twitch.Slowmode(c, userID, channelID, false)
		if err != nil {
			return err
		}

		return nil
	},
}

var slowmodeDurationCmd = &cobra.Command{
	Use:   "duration",
	Short: "Set slowmode with -c (channel name) and -d (duration in seconds) flags. Duration can be 3..120 seconds.",
	RunE: func(cmd *cobra.Command, args []string) error {
		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		duration, err := cmd.Flags().GetInt("duration")
		if err != nil {
			return err
		}

		if duration < 3 || duration > 120 {
			fmt.Printf("Duration in seconds can only be between 3 and 120.\n")
			return fmt.Errorf("use a valid duration")
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		userID, err := twitch.GetMyUserID(c)
		if err != nil {
			return err
		}

		err = twitch.SlowmodeDuration(c, userID, channelID, duration)
		if err != nil {
			return err
		}

		return nil
	},
}

var submodeCmd = &cobra.Command{
	Use:   "submode",
	Short: "Enable or Disable subscriber-only mode with -c (channel name) flag",
}

var submodeOnCmd = &cobra.Command{
	Use:   "on",
	Short: "Enable subcriber-only mode with -c (channel name) flag",
	RunE: func(cmd *cobra.Command, args []string) error {
		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		userID, err := twitch.GetMyUserID(c)
		if err != nil {
			return err
		}

		err = twitch.SubOnlyMode(c, userID, channelID, true)
		if err != nil {
			return err
		}

		return nil
	},
}

var submodeOffCmd = &cobra.Command{
	Use:   "off",
	Short: "Disable subcriber-only mode with -c (channel name) flag",
	RunE: func(cmd *cobra.Command, args []string) error {
		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelID, err := twitch.GetUserID(c, channelname)
		if err != nil {
			return err
		}

		userID, err := twitch.GetMyUserID(c)
		if err != nil {
			return err
		}

		err = twitch.SubOnlyMode(c, userID, channelID, false)
		if err != nil {
			return err
		}

		return nil
	},
}
