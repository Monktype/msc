package twitch

import (
	"fmt"

	"github.com/nicklaw5/helix/v2"
)

// CreatePoll creates a Twitch poll with the given title, duration, and options.
// Returns a poll ID and error.
func CreatePoll(c helix.Client, channelID string, title string, durationInSeconds int, options []string) (string, error) {
	// Convert options to a slice of PollChoiceParam
	var pollChoices []helix.PollChoiceParam
	for _, option := range options {
		if len(option) > 25 {
			fmt.Printf("Warning: Option '%s' exceeds the maximum length of 25 characters.\n", option)
			option = option[:25] // Truncate if necessary
		}
		pollChoices = append(pollChoices, helix.PollChoiceParam{Title: option})
	}

	poll, err := c.CreatePoll(&helix.CreatePollParams{
		BroadcasterID: channelID,
		Title:         title,
		Choices:       pollChoices,
		Duration:      durationInSeconds,
	})
	if err != nil {
		fmt.Printf("Creating a poll failed: %s\n", err)
		return "", err
	}
	if poll.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", poll)
		return "", fmt.Errorf("check status code information")
	}

	fmt.Printf("Poll created with ID: %s\n", poll.Data.Polls[0].ID)
	return poll.Data.Polls[0].ID, nil
}

// GetPolls gets polls from a channel ID.
func GetPolls(c helix.Client, channelID string) ([]helix.Poll, error) {
	var emptyPollResponse []helix.Poll

	polls, err := c.GetPolls(&helix.PollsParams{
		BroadcasterID: channelID,
	})
	if err != nil {
		fmt.Printf("Failed to get polls on channel ID %s: %s\n", channelID, err)
		return emptyPollResponse, err
	}
	if polls.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", polls)
		return emptyPollResponse, fmt.Errorf("check status code information")
	}

	return polls.Data.Polls, nil
}

// GetPoll gets a poll from a poll ID string and a channel ID string.
// Technically Twitch's API allows 1<=x<=20 polls to be requested in a single call, but
// the upstream library doesn't seem to implement their code in that way and I don't
// see an immediate need to request multiple specific polls in a single call right
// now, so I'm not going to try to change that upstream.
func GetPoll(c helix.Client, channelID string, pollID string) (helix.Poll, error) {
	var emptyPollResponse helix.Poll

	polls, err := c.GetPolls(&helix.PollsParams{
		BroadcasterID: channelID,
		ID:            pollID,
	})
	if err != nil {
		fmt.Printf("Failed to get poll on channel ID %s: %s\n", channelID, err)
		return emptyPollResponse, err
	}
	if polls.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", polls)
		return emptyPollResponse, fmt.Errorf("check status code information")
	}
	if len(polls.Data.Polls) == 0 {
		fmt.Printf("Poll %s was not found on channel %s.\n", pollID, channelID)
		return emptyPollResponse, nil
	}

	return polls.Data.Polls[0], nil
}

// EndPoll terminates a poll.
// Takes a Client, the string of the channel ID, and the string of the poll ID.
// Returns error.
func EndPoll(c helix.Client, channelID string, pollID string) error {
	resp, err := c.EndPoll(&helix.EndPollParams{
		BroadcasterID: channelID,
		ID:            pollID,
		Status:        "TERMINATED",
	})
	if err != nil {
		fmt.Printf("Failed to terminate poll: %s\n", err)
		return err
	}
	if resp.StatusCode >= 300 {
		fmt.Printf("Status code was bad: %v\n", resp)
		return fmt.Errorf("check status code information")
	}

	return nil
}
