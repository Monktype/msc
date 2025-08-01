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
