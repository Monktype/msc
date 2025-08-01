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
	_, err := c.SendChatAnnouncement(&helix.SendChatAnnouncementParams{
		BroadcasterID: channelID,
		ModeratorID:   userID,
		Color:         AnnouncementColorMap[color],
		Message:       message, // "Max 500 characters, truncated thereafter"
	})
	if err != nil {
		fmt.Printf("Announcement failed to send: %s\n", err)
		return err
	}

	return nil
}
