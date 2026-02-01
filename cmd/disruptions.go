package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"tfl/internal/display"
)

var disruptionsCmd = &cobra.Command{
	Use:     "disruptions",
	Aliases: []string{"delays"},
	Short:   "Show service disruptions",
	Long: `Display current service disruptions across the tube network.

Examples:
  tfl disruptions
  tfl delays
  tfl disruptions --format json`,
	Run: func(cmd *cobra.Command, args []string) {
		disruptions, err := client.GetDisruptions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if IsJSON() {
			display.PrintDisruptionsJSON(disruptions)
		} else {
			display.PrintDisruptions(disruptions)
		}
	},
}

func init() {
	rootCmd.AddCommand(disruptionsCmd)
}
