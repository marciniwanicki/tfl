package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"tfl/internal/display"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show tube line status",
	Long: `Display the current status of all London Underground and Elizabeth lines.

Examples:
  tfl status
  tfl status --format json`,
	Run: func(cmd *cobra.Command, args []string) {
		statuses, err := client.GetTubeStatus()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if IsJSON() {
			display.PrintLineStatusesJSON(statuses)
		} else {
			display.PrintLineStatuses(statuses)
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
