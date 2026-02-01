package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type checkResult struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check if API key is configured and valid",
	Long: `Verify that a TfL API key is set and can successfully authenticate with the API.

The API key can be provided via the TFL_APP_KEY environment variable or the --key flag.

Examples:
  tfl check
  tfl check --key YOUR_API_KEY
  tfl check --format json`,
	Run: func(cmd *cobra.Command, args []string) {
		if !client.HasKey() {
			if IsJSON() {
				printCheckResult(checkResult{Valid: false, Message: "No API key configured"})
			} else {
				fmt.Fprintln(os.Stderr, "No API key configured.")
				fmt.Fprintln(os.Stderr, "Set TFL_APP_KEY environment variable or use --key flag.")
			}
			os.Exit(1)
		}

		if !IsJSON() {
			fmt.Print("Validating API key... ")
		}

		if err := client.ValidateKey(); err != nil {
			if IsJSON() {
				printCheckResult(checkResult{Valid: false, Message: "API key validation failed", Error: err.Error()})
			} else {
				fmt.Fprintln(os.Stderr, "failed")
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
			os.Exit(1)
		}

		if IsJSON() {
			printCheckResult(checkResult{Valid: true, Message: "API key is valid"})
		} else {
			fmt.Println("valid")
		}
	},
}

func printCheckResult(result checkResult) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(result)
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
