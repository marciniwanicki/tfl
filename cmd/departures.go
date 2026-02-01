package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"tfl/internal/display"
	"tfl/internal/tfl"
)

var limit int
var match string

var departuresCmd = &cobra.Command{
	Use:   "departures <station> [line]",
	Short: "Show departures from a station",
	Long: `Show upcoming departures from a station.

Examples:
  tfl departures "liverpool street"
  tfl departures "liverpool street" elizabeth
  tfl departures paddington central
  tfl departures paddington -n 5
  tfl departures liverpool -m "elizabeth heathrow terminal 5"`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stationQuery := args[0]
		var lineFilter string
		if len(args) > 1 {
			lineFilter = strings.ToLower(args[1])
		}

		stops, err := client.SearchStopPoints(stationQuery)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error searching stations: %v\n", err)
			os.Exit(1)
		}

		if len(stops) == 0 {
			fmt.Fprintf(os.Stderr, "No stations found matching '%s'\n", stationQuery)
			os.Exit(1)
		}

		stop := selectBestMatch(stops, stationQuery)

		var arrivals []tfl.Arrival
		if lineFilter != "" && match == "" {
			arrivals, err = client.GetArrivals(stop.ID, lineFilter)
		} else {
			arrivals, err = client.GetAllArrivalsAtStop(stop.ID)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching arrivals: %v\n", err)
			os.Exit(1)
		}

		if match != "" {
			arrivals = filterByMatch(arrivals, match)
		}

		if lineFilter != "" && match != "" {
			arrivals = filterByLine(arrivals, lineFilter)
		}

		if limit > 0 && len(arrivals) > limit {
			arrivals = arrivals[:limit]
		}

		display.PrintArrivals(arrivals, stop.Name)
	},
}

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for stations",
	Long:  `Search for stations by name.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stops, err := client.SearchStopPoints(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		display.PrintStopPoints(stops)
	},
}

func filterByMatch(arrivals []tfl.Arrival, match string) []tfl.Arrival {
	words := strings.Fields(strings.ToLower(match))
	var filtered []tfl.Arrival
	for _, a := range arrivals {
		searchText := strings.ToLower(a.LineName + " " + a.DestinationName)
		allMatch := true
		for _, word := range words {
			if !strings.Contains(searchText, word) {
				allMatch = false
				break
			}
		}
		if allMatch {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

func filterByLine(arrivals []tfl.Arrival, line string) []tfl.Arrival {
	var filtered []tfl.Arrival
	for _, a := range arrivals {
		if strings.ToLower(a.LineName) == line {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

func selectBestMatch(stops []tfl.StopPoint, query string) tfl.StopPoint {
	query = strings.ToLower(query)

	for _, stop := range stops {
		if strings.ToLower(stop.Name) == query {
			return stop
		}
	}

	for _, stop := range stops {
		if strings.Contains(strings.ToLower(stop.Name), query) {
			return stop
		}
	}

	return stops[0]
}

func init() {
	departuresCmd.Flags().IntVarP(&limit, "limit", "n", 0, "Maximum number of departures to show")
	departuresCmd.Flags().StringVarP(&match, "match", "m", "", "Fuzzy filter by line name and/or destination")
	rootCmd.AddCommand(departuresCmd)
	rootCmd.AddCommand(searchCmd)
}
