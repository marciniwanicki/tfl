package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"tfl/internal/display"
	"tfl/internal/tfl"
)

var limit int
var match string
var departureTime string

var departuresCmd = &cobra.Command{
	Use:   "departures <station-name>",
	Short: "Show departures from a station",
	Long: `Show upcoming departures from a station.

Station names are matched case-insensitively and support partial matching.
Use quotes for station names containing spaces. Use -m to filter by line or destination.

Examples:
  tfl departures "Liverpool Street"
  tfl departures Paddington
  tfl departures Paddington -n 5
  tfl departures Paddington -m Central
  tfl departures Paddington -m "Heathrow Terminal 5"
  tfl departures Paddington --time 14:30
  tfl departures Paddington --format json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stationQuery := args[0]

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
		var minTime time.Time

		if departureTime != "" {
			minTime, err = parseTimeToday(departureTime)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing time: %v\n", err)
				os.Exit(1)
			}
		}

		// Use timetable if time is more than 30 minutes in the future
		useTimetable := departureTime != "" && time.Until(minTime) > 30*time.Minute
		timetableFailed := false

		if useTimetable {
			arrivals, err = getArrivalsFromTimetable(stop.ID, match, minTime)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching timetable: %v\n", err)
				os.Exit(1)
			}
			if len(arrivals) == 0 {
				timetableFailed = true
			}
		}

		// Fall back to real-time arrivals if timetable returned no results
		// (e.g., Elizabeth line doesn't support timetable API)
		if !useTimetable || timetableFailed {
			if timetableFailed {
				fmt.Println("Note: Timetable unavailable for this line. Real-time data only covers ~30 minutes ahead.")
			}

			arrivals, err = client.GetAllArrivalsAtStop(stop.ID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching arrivals: %v\n", err)
				os.Exit(1)
			}

			// Apply time filter - this may result in no arrivals if time is far in future
			if departureTime != "" {
				arrivals = filterByTime(arrivals, minTime)
			}
		}

		if match != "" {
			arrivals = filterByMatch(arrivals, match)
		}

		if limit > 0 && len(arrivals) > limit {
			arrivals = arrivals[:limit]
		}

		if IsJSON() {
			display.PrintArrivalsJSON(arrivals, stop.Name)
		} else {
			display.PrintArrivals(arrivals, stop.Name)
		}
	},
}

var searchCmd = &cobra.Command{
	Use:   "search <station-name>",
	Short: "Search for stations",
	Long: `Search for stations by name.

Station names are matched case-insensitively and support partial matching.

Examples:
  tfl search "King's Cross"
  tfl search Paddington
  tfl search Victoria
  tfl search Liverpool --format json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stops, err := client.SearchStopPoints(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		if IsJSON() {
			display.PrintStopPointsJSON(stops)
		} else {
			display.PrintStopPoints(stops)
		}
	},
}

func filterByMatch(arrivals []tfl.Arrival, match string) []tfl.Arrival {
	words := strings.Fields(strings.ToLower(match))
	var filtered []tfl.Arrival
	for _, a := range arrivals {
		searchText := strings.ToLower(a.LineName + " " + a.DestinationName + " " + a.PlatformName)
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

func parseTimeToday(timeStr string) (time.Time, error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format, use HH:MM (e.g., 14:30)")
	}
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location()), nil
}

func filterByTime(arrivals []tfl.Arrival, minTime time.Time) []tfl.Arrival {
	var filtered []tfl.Arrival
	for _, a := range arrivals {
		if !a.ExpectedArrival.Before(minTime) {
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

func getArrivalsFromTimetable(stopID, lineFilter string, minTime time.Time) ([]tfl.Arrival, error) {
	detail, err := client.GetStopPointDetails(stopID)
	if err != nil {
		return nil, err
	}

	// Collect all line+stop combinations to fetch
	type lineStop struct {
		lineID   string
		lineName string
		stopID   string
	}
	var lineStops []lineStop

	// Check children for tube/elizabeth-line stops
	for _, child := range detail.Children {
		for _, line := range child.Lines {
			if lineFilter == "" || strings.ToLower(line.ID) == lineFilter || strings.Contains(strings.ToLower(line.Name), lineFilter) {
				lineStops = append(lineStops, lineStop{
					lineID:   line.ID,
					lineName: line.Name,
					stopID:   child.ID,
				})
			}
		}
	}

	// Also check the stop itself if it has lines
	for _, line := range detail.Lines {
		if lineFilter == "" || strings.ToLower(line.ID) == lineFilter || strings.Contains(strings.ToLower(line.Name), lineFilter) {
			lineStops = append(lineStops, lineStop{
				lineID:   line.ID,
				lineName: line.Name,
				stopID:   detail.ID,
			})
		}
	}

	var allArrivals []tfl.Arrival
	seen := make(map[string]bool)

	// First pass: fetch all timetables and collect station names
	stationNames := make(map[string]string)
	var timetables []*tfl.TimetableResponse

	for _, ls := range lineStops {
		if seen[ls.lineID] {
			continue
		}
		seen[ls.lineID] = true

		// Fetch both directions
		for _, direction := range []string{"inbound", "outbound"} {
			timetable, err := client.GetTimetable(ls.lineID, ls.stopID, direction)
			if err != nil {
				continue
			}

			// Collect station names
			for _, s := range timetable.Stations {
				name := s.Name
				name = strings.TrimSuffix(name, " Underground Station")
				name = strings.TrimSuffix(name, " Rail Station")
				name = strings.TrimSuffix(name, " DLR Station")
				stationNames[s.ID] = name
			}

			timetables = append(timetables, timetable)
		}
	}

	// Second pass: parse all timetables with complete station names
	for _, timetable := range timetables {
		arrivals := parseTimetableWithStations(timetable, minTime, stationNames)
		allArrivals = append(allArrivals, arrivals...)
	}

	// Sort by expected arrival time
	sort.Slice(allArrivals, func(i, j int) bool {
		return allArrivals[i].ExpectedArrival.Before(allArrivals[j].ExpectedArrival)
	})

	return allArrivals, nil
}

func parseTimetableWithStations(tt *tfl.TimetableResponse, minTime time.Time, stationNames map[string]string) []tfl.Arrival {
	if tt == nil || len(tt.Timetable.Routes) == 0 {
		return nil
	}

	now := time.Now()
	today := now.Weekday()
	var arrivals []tfl.Arrival

	for _, route := range tt.Timetable.Routes {
		// Build destination map from station intervals
		destMap := make(map[int]string)
		for _, si := range route.StationIntervals {
			id, _ := strconv.Atoi(si.ID)
			if len(si.Intervals) > 0 {
				// Last station in the interval is the destination
				lastStopID := si.Intervals[len(si.Intervals)-1].StopID
				if name, ok := stationNames[lastStopID]; ok {
					destMap[id] = name
				} else {
					destMap[id] = formatStopID(lastStopID)
				}
			}
		}

		for _, schedule := range route.Schedules {
			if !scheduleMatchesDay(schedule.Name, today) {
				continue
			}

			for _, journey := range schedule.KnownJourneys {
				hour, _ := strconv.Atoi(journey.Hour)
				minute, _ := strconv.Atoi(journey.Minute)

				departTime := time.Date(now.Year(), now.Month(), now.Day(),
					hour, minute, 0, 0, now.Location())

				// Skip departures before the requested time
				if departTime.Before(minTime) {
					continue
				}

				// Skip departures more than 4 hours after the requested time
				if departTime.After(minTime.Add(4 * time.Hour)) {
					continue
				}

				dest := destMap[journey.IntervalID]
				if dest == "" {
					dest = tt.Direction
				}

				arrivals = append(arrivals, tfl.Arrival{
					LineName:        tt.LineName,
					LineID:          tt.LineID,
					DestinationName: dest,
					ExpectedArrival: departTime,
					TimeToStation:   int(time.Until(departTime).Seconds()),
				})
			}
		}
	}

	return arrivals
}

func scheduleMatchesDay(scheduleName string, day time.Weekday) bool {
	name := strings.ToLower(scheduleName)

	switch day {
	case time.Saturday:
		return strings.Contains(name, "saturday")
	case time.Sunday:
		return strings.Contains(name, "sunday")
	case time.Friday:
		if strings.Contains(name, "friday") {
			return true
		}
		return strings.Contains(name, "monday") && strings.Contains(name, "friday")
	default: // Monday-Thursday
		return strings.Contains(name, "monday")
	}
}

func formatStopID(stopID string) string {
	// Common station code mappings for cases where API IDs don't match
	codeMap := map[string]string{
		"940GZZLUWRP": "West Ruislip",
		"940GZZLUEBY": "Ealing Broadway",
		"940GZZLUEAN": "East Acton",
		"940GZZLUNOA": "North Acton",
		"940GZZLUWCY": "White City",
		"940GZZLUHLT": "Hainault",
		"940GZZLUEPG": "Epping",
	}

	if name, ok := codeMap[stopID]; ok {
		return name
	}

	// Fallback: clean up the ID
	stopID = strings.TrimPrefix(stopID, "940GZZLU")
	stopID = strings.TrimPrefix(stopID, "910G")
	return stopID
}

func init() {
	departuresCmd.Flags().IntVarP(&limit, "limit", "n", 0, "Maximum number of departures to show")
	departuresCmd.Flags().StringVarP(&match, "match", "m", "", "Fuzzy filter by line name and/or destination")
	departuresCmd.Flags().StringVarP(&departureTime, "time", "t", "", "Show departures at or after this time (HH:MM)")
	rootCmd.AddCommand(departuresCmd)
	rootCmd.AddCommand(searchCmd)
}
