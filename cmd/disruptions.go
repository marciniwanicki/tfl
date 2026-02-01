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
	Long:    `Display current service disruptions across the tube network.`,
	Run: func(cmd *cobra.Command, args []string) {
		disruptions, err := client.GetDisruptions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		display.PrintDisruptions(disruptions)
	},
}

func init() {
	rootCmd.AddCommand(disruptionsCmd)
}
