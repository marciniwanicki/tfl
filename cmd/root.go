package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"tfl/internal/tfl"
)

var (
	appKey string
	client *tfl.Client
)

var rootCmd = &cobra.Command{
	Use:   "tfl",
	Short: "Transport for London CLI",
	Long:  `A command-line interface for Transport for London services.`,
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
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}
