package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if API key is configured and valid",
	Long:  `Verify that a TfL API key is set and can successfully authenticate with the API.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !client.HasKey() {
			fmt.Fprintln(os.Stderr, "No API key configured.")
			fmt.Fprintln(os.Stderr, "Set TFL_APP_KEY environment variable or use --key flag.")
			os.Exit(1)
		}

		fmt.Print("Validating API key... ")
		if err := client.ValidateKey(); err != nil {
			fmt.Fprintln(os.Stderr, "failed")
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("valid")
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
