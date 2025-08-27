package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/monktype/msc/twitch"
	"github.com/nicklaw5/helix/v2"
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
		if len(args) > 5 {
			fmt.Printf("At most 5 poll options are allowed, %d provided.\n", len(args))
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
		resultstring, err := watchPollCompletion(c, userID, pollID)
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

// --- Code here is for watching for responses in the CLI ---

type WatchPollResult struct {
	Result string
	Error  error
}

// Internal watchPollCompletion worker function (for terminating the poll)
// Re-using watchPollResult as a result, but only the Error component is going to be used.
func watchPollCompletionTerminationWorker(c helix.Client, channelID string, pollID string, resultChan chan<- WatchPollResult, doneChan <-chan os.Signal) {
	defer close(resultChan)
	for {
		select {
		case <-doneChan:
			fmt.Printf("Terminating poll...\n")
			err := twitch.EndPoll(c, channelID, pollID)
			resultChan <- WatchPollResult{Result: "", Error: err}
			return
		}
	}
}

// Internal watchPollCompletion worker function
// NOTE: It could be broken into some smaller functions if desired, but not critical right now.
func watchPollCompletionWorker(c helix.Client, channelID string, pollID string, resultChan chan<- WatchPollResult) {
	defer close(resultChan)
	pollGetFailCount := 0
	for {
		// Fetch the poll status
		poll, err := twitch.GetPoll(c, channelID, pollID)
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

		if poll.ID != "" {
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
		} else /* this triggers if `poll.ID == ""` */ {
			fmt.Printf("Poll %s not found yet...\n", pollID)
		}
		// Wait for a second before checking again
		time.Sleep(1 * time.Second)
	}
}

// watchPollCompletion checks the status of the poll and prints the results when completed.
func watchPollCompletion(c helix.Client, channelID string, pollID string) (string, error) {
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
