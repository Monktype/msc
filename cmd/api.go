package cmd

import (
	"fmt"

	"github.com/monktype/msc/api"
	"github.com/spf13/cobra"
)

var startApiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start API",
	RunE: func(cmd *cobra.Command, args []string) error {
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			return err
		}
		// Not checking for it to be a valid port, the computer will error for me instead of me spending time to type a checker (I typed this comment in the time I saved)

		fmt.Printf("Starting API server on port %d...\n", port)

		err = api.ApiServer(port)
		if err != nil {
			return err
		}

		return nil
	},
}
