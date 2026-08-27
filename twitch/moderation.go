package twitch

import (
	"fmt"

	"github.com/nicklaw5/helix/v2"
)

// BannedUsersPage is one page of banned/timed-out users returned by GetBannedUsers.
// Cursor is the cursor to pass as `after` to GetBannedUsers to fetch the next page;
// it is "" when there is no next page.
type BannedUsersPage struct {
	Bans   []helix.Ban `json:"bans"`
	Cursor string      `json:"cursor"`
}

// GetBannedUsers gets the banned and timed-out users for the given channel ID.
// `after` is an optional pagination cursor to continue from a previous page; the
// returned page's Cursor carries the cursor for the next page ("" when done).
// Notes for downstream tools:
//   - The returned list includes both perma-bans and timeouts; a zero `ExpiresAt`
//     means a perma-ban.
//   - This endpoint has no `first` param; the page size is the API default (100).
func GetBannedUsers(c *helix.Client, channelID string, after string) (BannedUsersPage, error) {
	var page BannedUsersPage

	resp, err := c.GetBannedUsers(&helix.BannedUsersParams{
		BroadcasterID: channelID,
		After:         after,
	})
	if err != nil {
		fmt.Printf("Getting banned users failed: %s\n", err)
		return page, err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return page, fmt.Errorf("check status code information")
	}

	page.Bans = resp.Data.Bans
	page.Cursor = resp.Data.Pagination.Cursor
	return page, nil
}

// GetAllBannedUsers returns every banned/timed-out user on the channel by fetching
// all pages. It is the convenience wrapper around GetBannedUsers for callers that
// don't want to manage the cursor loop themselves. It keeps going until the cursor
// is exhausted; if the cursor ever stops advancing (the API returns a cursor it
// already gave), it stops and returns the accumulated list along with a warning
// rather than looping forever on the same page.
func GetAllBannedUsers(c *helix.Client, channelID string) ([]helix.Ban, error) {
	var all []helix.Ban
	after := ""
	seen := map[string]bool{} // cursors already received; a repeat means we're looping
	for {
		p, err := GetBannedUsers(c, channelID, after)
		if err != nil {
			return nil, err
		}
		all = append(all, p.Bans...)
		if p.Cursor == "" {
			break // no next page — done
		}
		if seen[p.Cursor] {
			fmt.Printf("Warning: GetAllBannedUsers stopped early; the cursor stopped advancing and the list may be incomplete.\n")
			break
		}
		seen[p.Cursor] = true
		after = p.Cursor
	}
	return all, nil
}

// BanUser bans or times out a user in the given channel.
// `duration` is the timeout length in seconds; 0 means a perma-ban.
// `reason` is the optional ban reason.
func BanUser(c *helix.Client, moderatorID string, channelID string, targetID string, duration int, reason string) error {
	resp, err := c.BanUser(&helix.BanUserParams{
		BroadcasterID: channelID,
		ModeratorId:   moderatorID,
		Body: helix.BanUserRequestBody{
			Duration: duration,
			Reason:   reason,
			UserId:   targetID,
		},
	})
	if err != nil {
		fmt.Printf("Banning user failed: %s\n", err)
		return err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return fmt.Errorf("check status code information")
	}

	return nil
}

// UnbanUser removes the ban or timeout from the given user in the given channel.
// Note: there is no batch-unban in the Twitch API, so removing many users requires
// one call per user, which is rate-limited.
func UnbanUser(c *helix.Client, moderatorID string, channelID string, targetID string) error {
	resp, err := c.UnbanUser(&helix.UnbanUserParams{
		BroadcasterID: channelID,
		ModeratorID:   moderatorID,
		UserID:        targetID,
	})
	if err != nil {
		fmt.Printf("Unbanning user failed: %s\n", err)
		return err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return fmt.Errorf("check status code information")
	}

	return nil
}
