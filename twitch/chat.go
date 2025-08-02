package twitch

import (
	"fmt"

	"github.com/nicklaw5/helix/v2"
)

type AnnouncementColor int

const (
	AnnouncementColorPrimary AnnouncementColor = iota
	AnnouncementColorBlue
	AnnouncementColorGreen
	AnnouncementColorOrange
	AnnouncementColorPurple
)

var AnnouncementColorMap = map[AnnouncementColor]string{
	AnnouncementColorPrimary: "primary", // "channel accent color"
	AnnouncementColorBlue:    "blue",
	AnnouncementColorGreen:   "green",
	AnnouncementColorOrange:  "orange",
	AnnouncementColorPurple:  "purple",
}

func SendAnnouncement(c helix.Client, userID string, channelID string, color AnnouncementColor, message string) error {
	resp, err := c.SendChatAnnouncement(&helix.SendChatAnnouncementParams{
		BroadcasterID: channelID,
		ModeratorID:   userID,
		Color:         AnnouncementColorMap[color],
		Message:       message, // "Max 500 characters, truncated thereafter"
	})
	if err != nil {
		fmt.Printf("Announcement failed to send: %s\n", err)
		return err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return fmt.Errorf("check status code information")
	}

	return nil
}

func SendShoutout(c helix.Client, userID string, channelID string, targetID string) error {
	resp, err := c.SendShoutout(&helix.SendShoutoutParams{
		FromBroadcasterID: channelID,
		ToBroadcasterID:   targetID,
		ModeratorID:       userID,
	})
	if err != nil {
		fmt.Printf("Shoutout failed to send: %s\n", err)
		return err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return fmt.Errorf("check status code information")
	}

	return nil
}

func EmoteOnly(c helix.Client, userID string, channelID string, state bool) error {
	resp, err := c.UpdateChatSettings(&helix.UpdateChatSettingsParams{
		ModeratorID:   userID,
		BroadcasterID: channelID,
		EmoteMode:     &state,
	})
	if err != nil {
		var statestring string
		if state {
			statestring = "on"
		} else {
			statestring = "off"
		}
		fmt.Printf("Set emote mode \"%s\" failed: %s\n", statestring, err)
		return err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return fmt.Errorf("check status code information")
	}

	return nil
}

func FollowerOnly(c helix.Client, userID string, channelID string, state bool) error {
	resp, err := c.UpdateChatSettings(&helix.UpdateChatSettingsParams{
		ModeratorID:   userID,
		BroadcasterID: channelID,
		FollowerMode:  &state,
	})
	if err != nil {
		var statestring string
		if state {
			statestring = "on"
		} else {
			statestring = "off"
		}
		fmt.Printf("Set follower mode \"%s\" failed: %s\n", statestring, err)
		return err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return fmt.Errorf("check status code information")
	}

	return nil
}

func FollowerOnlyDuration(c helix.Client, userID string, channelID string, duration int) error {
	trueFlagBecauseItWantsAVariable := true
	resp, err := c.UpdateChatSettings(&helix.UpdateChatSettingsParams{
		ModeratorID:          userID,
		BroadcasterID:        channelID,
		FollowerMode:         &trueFlagBecauseItWantsAVariable,
		FollowerModeDuration: &duration,
	})
	if err != nil {
		fmt.Printf("Set follower mode \"on\" to duration %d failed: %s\n", duration, err)
		return err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return fmt.Errorf("check status code information")
	}

	return nil
}

func Slowmode(c helix.Client, userID string, channelID string, state bool) error {
	resp, err := c.UpdateChatSettings(&helix.UpdateChatSettingsParams{
		ModeratorID:   userID,
		BroadcasterID: channelID,
		SlowMode:      &state,
	})
	if err != nil {
		var statestring string
		if state {
			statestring = "on"
		} else {
			statestring = "off"
		}
		fmt.Printf("Set slowmode \"%s\" failed: %s\n", statestring, err)
		return err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return fmt.Errorf("check status code information")
	}

	return nil
}

func SlowmodeDuration(c helix.Client, userID string, channelID string, duration int) error {
	trueFlagBecauseItWantsAVariable := true
	resp, err := c.UpdateChatSettings(&helix.UpdateChatSettingsParams{
		ModeratorID:      userID,
		BroadcasterID:    channelID,
		SlowMode:         &trueFlagBecauseItWantsAVariable,
		SlowModeWaitTime: &duration,
	})
	if err != nil {
		fmt.Printf("Set slowmode \"on\" to duration %d failed: %s\n", duration, err)
		return err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return fmt.Errorf("check status code information")
	}

	return nil
}

func SubOnlyMode(c helix.Client, userID string, channelID string, state bool) error {
	resp, err := c.UpdateChatSettings(&helix.UpdateChatSettingsParams{
		ModeratorID:    userID,
		BroadcasterID:  channelID,
		SubscriberMode: &state,
	})
	if err != nil {
		var statestring string
		if state {
			statestring = "on"
		} else {
			statestring = "off"
		}
		fmt.Printf("Set subscriber-only mode \"%s\" failed: %s\n", statestring, err)
		return err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return fmt.Errorf("check status code information")
	}

	return nil
}
