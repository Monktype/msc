package cmd

import (
	"fmt"

	"github.com/monktype/msc/twitch"
	"github.com/nicklaw5/helix/v2"
	"github.com/spf13/cobra"
)

var startadCmd = &cobra.Command{
	Use:   "start-ad",
	Short: "Start Advertisements",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := twitch.GetClient()
		if err != nil {
			return err
		}

		channelname, err := cmd.Flags().GetString("channel-name")
		if err != nil {
			return err
		}

		length, err := cmd.Flags().GetInt("length")
		if err != nil {
			return err
		}

		var lengthEnum helix.AdLengthEnum

		switch length {
		case 30:
			lengthEnum = helix.AdLen30
		case 60:
			lengthEnum = helix.AdLen60
		case 90:
			lengthEnum = helix.AdLen90
		case 120:
			lengthEnum = helix.AdLen120
		case 150:
			lengthEnum = helix.AdLen150
		case 180:
			lengthEnum = helix.AdLen180
		default:
			fmt.Printf("Length %d is invalid; only 30, 60, 90, 120, 150, 180 are valid values.\n", length)
			return fmt.Errorf("use correct length value")
		}

		err = twitch.StartCommercial(c, channelname, lengthEnum)
		if err != nil {
			return err
		}

		return nil
	},
}
