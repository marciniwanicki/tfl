package display

import (
	"encoding/json"
	"fmt"
	"os"

	"tfl/internal/tfl"
)

type ArrivalJSON struct {
	Line            string `json:"line"`
	LineID          string `json:"line_id"`
	Destination     string `json:"destination"`
	Platform        string `json:"platform,omitempty"`
	TimeToStation   int    `json:"time_to_station_seconds"`
	MinutesAway     int    `json:"minutes_away"`
	ExpectedArrival string `json:"expected_arrival"`
}

type DeparturesOutput struct {
	Station   string        `json:"station"`
	Arrivals  []ArrivalJSON `json:"arrivals"`
	Count     int           `json:"count"`
}

type LineStatusJSON struct {
	Line        string `json:"line"`
	LineID      string `json:"line_id"`
	Status      string `json:"status"`
	Severity    int    `json:"severity"`
	Reason      string `json:"reason,omitempty"`
}

type StatusOutput struct {
	Lines []LineStatusJSON `json:"lines"`
	Count int              `json:"count"`
}

type DisruptionJSON struct {
	Category    string `json:"category"`
	Description string `json:"description"`
}

type DisruptionsOutput struct {
	Disruptions []DisruptionJSON `json:"disruptions"`
	Count       int              `json:"count"`
}

type StopPointJSON struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Zone  string   `json:"zone,omitempty"`
	Modes []string `json:"modes"`
}

type StopPointsOutput struct {
	Stations []StopPointJSON `json:"stations"`
	Count    int             `json:"count"`
}

func printJSON(v interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
	}
}

func PrintArrivalsJSON(arrivals []tfl.Arrival, stationName string) {
	output := DeparturesOutput{
		Station:  stationName,
		Arrivals: make([]ArrivalJSON, 0, len(arrivals)),
		Count:    len(arrivals),
	}

	for _, arr := range arrivals {
		output.Arrivals = append(output.Arrivals, ArrivalJSON{
			Line:            arr.LineName,
			LineID:          arr.LineID,
			Destination:     arr.DestinationName,
			Platform:        arr.PlatformName,
			TimeToStation:   arr.TimeToStation,
			MinutesAway:     arr.TimeToStation / 60,
			ExpectedArrival: arr.ExpectedArrival.Local().Format("15:04"),
		})
	}

	printJSON(output)
}

func PrintLineStatusesJSON(statuses []tfl.LineStatus) {
	output := StatusOutput{
		Lines: make([]LineStatusJSON, 0, len(statuses)),
		Count: len(statuses),
	}

	for _, line := range statuses {
		status := line.LineStatuses[0]
		output.Lines = append(output.Lines, LineStatusJSON{
			Line:     line.Name,
			LineID:   line.ID,
			Status:   status.StatusSeverityDescription,
			Severity: status.StatusSeverity,
			Reason:   status.Reason,
		})
	}

	printJSON(output)
}

func PrintDisruptionsJSON(disruptions []tfl.Disruption) {
	output := DisruptionsOutput{
		Disruptions: make([]DisruptionJSON, 0, len(disruptions)),
		Count:       len(disruptions),
	}

	for _, d := range disruptions {
		output.Disruptions = append(output.Disruptions, DisruptionJSON{
			Category:    d.CategoryDescription,
			Description: d.Description,
		})
	}

	printJSON(output)
}

func PrintStopPointsJSON(stops []tfl.StopPoint) {
	output := StopPointsOutput{
		Stations: make([]StopPointJSON, 0, len(stops)),
		Count:    len(stops),
	}

	for _, stop := range stops {
		output.Stations = append(output.Stations, StopPointJSON{
			ID:    stop.ID,
			Name:  stop.Name,
			Zone:  stop.Zone,
			Modes: stop.Modes,
		})
	}

	printJSON(output)
}
