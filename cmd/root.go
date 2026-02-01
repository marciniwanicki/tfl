package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"tfl/internal/tfl"
)

var (
	appKey       string
	client       *tfl.Client
	outputFormat string
)

var rootCmd = &cobra.Command{
	Use:   "tfl",
	Short: "Transport for London CLI",
	Long: `A command-line interface for Transport for London services.

Station names are matched case-insensitively and support partial matching,
so "paddington", "Paddington", and "padd" will all find Paddington station.

Examples:
  tfl status                              Show all tube line statuses
  tfl disruptions                         Show current service disruptions
  tfl departures "Liverpool Street"       Show departures from a station
  tfl departures Paddington Central -n 5  Filter by line, limit results
  tfl search "King's Cross"               Search for stations`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if appKey == "" {
			appKey = os.Getenv("TFL_APP_KEY")
		}
		client = tfl.NewClient(appKey)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&appKey, "key", "", "TfL API key (or set TFL_APP_KEY env var)")
	rootCmd.PersistentFlags().StringVar(&outputFormat, "format", "text", "Output format: text or json")
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}

func IsJSON() bool {
	return outputFormat == "json"
}
