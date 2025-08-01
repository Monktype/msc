package twitch

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	return polls.Data.Polls, nil
}

type WatchPollResult struct {
	Result string
	Error  error
}

// Internal WatchPollCompletion worker function (for terminating the poll)
// Re-using WatchPollResult as a result, but only the Error component is going to be used.
func watchPollCompletionTerminationWorker(c helix.Client, channelID string, pollID string, resultChan chan<- WatchPollResult, doneChan <-chan os.Signal) {
	defer close(resultChan)
	for {
		select {
		case <-doneChan:
			fmt.Printf("Terminating poll...\n")
			err := EndPoll(c, channelID, pollID)
			resultChan <- WatchPollResult{Result: "", Error: err}
			return
		}
	}
}

// Internal WatchPollCompletion worker function
func watchPollCompletionWorker(c helix.Client, channelID string, pollID string, resultChan chan<- WatchPollResult) {
	defer close(resultChan)
	pollGetFailCount := 0
	for {
		// Fetch the poll status
		polls, err := GetPolls(c, channelID)
		if err != nil {
			fmt.Printf("Failed getting polls on channel ID %s: %s\n", channelID, err)
			pollGetFailCount = pollGetFailCount + 1
			if pollGetFailCount > 2 {
				resultChan <- WatchPollResult{Result: "", Error: err}
				return
			}
			fmt.Printf("Trying again...\n")
			time.Sleep(1 * time.Second)
			continue
		}

		pollfound := false

		for _, poll := range polls {
			if poll.ID == pollID {
				pollfound = true
				// Check if the poll is completed
				if poll.Status != "ACTIVE" { // There are many statuses that mean "not running," but "ACTIVE" is "running"
					fmt.Printf("Poll completed! Here are the results for \"%s\":\n", poll.Title)

					maxVotes := 0
					var winningOptions []string

					for _, option := range poll.Choices {
						fmt.Printf("Option: %s, Votes: %d\n", option.Title, option.Votes)
						if option.Votes > maxVotes {
							maxVotes = option.Votes
							winningOptions = []string{option.Title}
						} else if option.Votes == maxVotes {
							winningOptions = append(winningOptions, option.Title)
						}
					}

					fmt.Printf("\n")

					var resultstring string

					if len(winningOptions) == 1 {
						resultstring = fmt.Sprintf("The winning option (at %d votes) is: %s", maxVotes, winningOptions[0])
					} else {
						resultstring = fmt.Sprintf("The top tie options (at %d votes) are:", maxVotes)
						for i := range winningOptions {
							if i == 0 {
								resultstring = resultstring + fmt.Sprintf(" %s", winningOptions[i])
							} else {
								resultstring = resultstring + fmt.Sprintf("; %s", winningOptions[i])
							}
						}
					}

					resultChan <- WatchPollResult{Result: resultstring, Error: nil}
					return
				}
				break
			}
		}
		if !pollfound {
			fmt.Printf("Poll %s not found yet...\n", pollID)
		}
		// Wait for a second before checking again
		time.Sleep(1 * time.Second)
	}
}

// WatchPollCompletion checks the status of the poll and prints the results when completed.
func WatchPollCompletion(c helix.Client, channelID string, pollID string) (string, error) {
	resultChan := make(chan WatchPollResult)
	termResultChan := make(chan WatchPollResult)
	doneChan := make(chan os.Signal, 1)
	signal.Notify(doneChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the termination worker
	// It catches a ctrl+c and sends it over doneChan.
	// That then sends a terminate to Twitch, which will then picked up by the other goroutine
	// and cause it to naturally close / look for results / etc.
	go watchPollCompletionTerminationWorker(c, channelID, pollID, termResultChan, doneChan)

	// Start the regular worker
	// It gathers results and returns them as soon as the poll is not "ACTIVE" anymore.
	go watchPollCompletionWorker(c, channelID, pollID, resultChan)

	fmt.Printf("Press CTRL+C once to close the poll and tally results early.\n\n")

	for {
		select {
		case result := <-resultChan:
			close(doneChan)
			if result.Error != nil {
				return "", result.Error
			}
			return result.Result, nil
		case termResult := <-termResultChan:
			if termResult.Error != nil {
				fmt.Printf("Failed to terminate. If you see this, CTRL+C more to terminate the program or wait for the poll to finish. %s\n", termResult.Error)
			}
		}
	}
}

// EndPoll terminates a poll.
// Takes a Client, the string of the channel ID, and the string of the poll ID.
// Returns error.
func EndPoll(c helix.Client, channelID string, pollID string) error {
	_, err := c.EndPoll(&helix.EndPollParams{
		BroadcasterID: channelID,
		ID:            pollID,
		Status:        "TERMINATED",
	})
	if err != nil {
		fmt.Printf("Failed to terminate poll: %s\n", err)
		return err
	}

	return nil
}
