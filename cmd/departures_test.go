package cmd

import (
	"testing"
	"time"

	"tfl/internal/tfl"
)

func TestFilterByMatch(t *testing.T) {
	arrivals := []tfl.Arrival{
		{LineName: "Elizabeth", DestinationName: "Heathrow Terminal 5", PlatformName: "Platform 1"},
		{LineName: "Elizabeth", DestinationName: "Shenfield", PlatformName: "Platform 2"},
		{LineName: "Central", DestinationName: "Ealing Broadway", PlatformName: "Westbound"},
		{LineName: "District", DestinationName: "Richmond", PlatformName: "Platform 3"},
	}

	tests := []struct {
		name     string
		match    string
		expected int
	}{
		{"single word matches line", "elizabeth", 2},
		{"single word matches destination", "heathrow", 1},
		{"multiple words same field", "heathrow terminal", 1},
		{"multiple words across fields", "elizabeth heathrow", 1},
		{"case insensitive", "ELIZABETH HEATHROW", 1},
		{"partial match", "eliz heat", 1},
		{"no match", "northern", 0},
		{"all words must match", "elizabeth richmond", 0},
		{"empty match returns all", "", 4},
		{"matches platform", "westbound", 1},
		{"matches line and platform", "elizabeth platform 1", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterByMatch(arrivals, tt.match)
			if len(result) != tt.expected {
				t.Errorf("filterByMatch(%q) = %d arrivals, want %d", tt.match, len(result), tt.expected)
			}
		})
	}
}

func TestFilterByLine(t *testing.T) {
	arrivals := []tfl.Arrival{
		{LineName: "Elizabeth", DestinationName: "Heathrow Terminal 5"},
		{LineName: "Elizabeth", DestinationName: "Shenfield"},
		{LineName: "Central", DestinationName: "Ealing Broadway"},
		{LineName: "District", DestinationName: "Richmond"},
	}

	tests := []struct {
		name     string
		line     string
		expected int
	}{
		{"matches elizabeth", "elizabeth", 2},
		{"matches central", "central", 1},
		{"no match", "northern", 0},
		{"case sensitive input assumed lowercase", "Elizabeth", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterByLine(arrivals, tt.line)
			if len(result) != tt.expected {
				t.Errorf("filterByLine(%q) = %d arrivals, want %d", tt.line, len(result), tt.expected)
			}
		})
	}
}

func TestParseTimeToday(t *testing.T) {
	tests := []struct {
		name      string
		timeStr   string
		wantHour  int
		wantMin   int
		wantError bool
	}{
		{"valid time", "14:30", 14, 30, false},
		{"midnight", "00:00", 0, 0, false},
		{"end of day", "23:59", 23, 59, false},
		{"single digit hour", "2:30", 2, 30, false},
		{"invalid format no colon", "1430", 0, 0, true},
		{"invalid hour", "25:00", 0, 0, true},
		{"invalid minute", "14:60", 0, 0, true},
		{"empty string", "", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseTimeToday(tt.timeStr)
			if tt.wantError {
				if err == nil {
					t.Errorf("parseTimeToday(%q) expected error, got nil", tt.timeStr)
				}
				return
			}
			if err != nil {
				t.Errorf("parseTimeToday(%q) unexpected error: %v", tt.timeStr, err)
				return
			}
			if result.Hour() != tt.wantHour || result.Minute() != tt.wantMin {
				t.Errorf("parseTimeToday(%q) = %02d:%02d, want %02d:%02d",
					tt.timeStr, result.Hour(), result.Minute(), tt.wantHour, tt.wantMin)
			}
			// Verify it's today's date
			now := time.Now()
			if result.Year() != now.Year() || result.Month() != now.Month() || result.Day() != now.Day() {
				t.Errorf("parseTimeToday(%q) date is not today", tt.timeStr)
			}
		})
	}
}

func TestFilterByTime(t *testing.T) {
	now := time.Now()
	baseTime := time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, now.Location())

	arrivals := []tfl.Arrival{
		{LineName: "Elizabeth", ExpectedArrival: baseTime.Add(-30 * time.Minute)}, // 13:30
		{LineName: "Central", ExpectedArrival: baseTime},                          // 14:00
		{LineName: "District", ExpectedArrival: baseTime.Add(30 * time.Minute)},   // 14:30
		{LineName: "Northern", ExpectedArrival: baseTime.Add(60 * time.Minute)},   // 15:00
	}

	tests := []struct {
		name     string
		minTime  time.Time
		expected int
	}{
		{"before all", baseTime.Add(-60 * time.Minute), 4},
		{"at first excluded", baseTime.Add(-30 * time.Minute), 4},
		{"between first and second", baseTime.Add(-15 * time.Minute), 3},
		{"at second", baseTime, 3},
		{"at third", baseTime.Add(30 * time.Minute), 2},
		{"at last", baseTime.Add(60 * time.Minute), 1},
		{"after all", baseTime.Add(90 * time.Minute), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterByTime(arrivals, tt.minTime)
			if len(result) != tt.expected {
				t.Errorf("filterByTime() = %d arrivals, want %d", len(result), tt.expected)
			}
		})
	}
}

func TestSelectBestMatch(t *testing.T) {
	stops := []tfl.StopPoint{
		{ID: "1", Name: "Liverpool Street Underground Station"},
		{ID: "2", Name: "Liverpool"},
		{ID: "3", Name: "Liverpool Street"},
	}

	tests := []struct {
		name   string
		query  string
		wantID string
	}{
		{"exact match", "Liverpool", "2"},
		{"exact match case insensitive", "liverpool", "2"},
		{"contains match", "Liverpool Street", "3"},
		{"no exact returns first contains", "Underground", "1"},
		{"no match returns first", "Paddington", "1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := selectBestMatch(stops, tt.query)
			if result.ID != tt.wantID {
				t.Errorf("selectBestMatch(%q) = %s, want %s", tt.query, result.ID, tt.wantID)
			}
		})
	}
}

func TestScheduleMatchesDay(t *testing.T) {
	tests := []struct {
		name         string
		scheduleName string
		day          time.Weekday
		want         bool
	}{
		{"saturday matches saturday", "Saturday (also Good Friday)", time.Saturday, true},
		{"saturday doesnt match sunday", "Saturday (also Good Friday)", time.Sunday, false},
		{"sunday matches sunday", "Sunday", time.Sunday, true},
		{"sunday doesnt match monday", "Sunday", time.Monday, false},
		{"friday matches friday", "Friday", time.Friday, true},
		{"monday-thursday matches monday", "Monday - Thursday", time.Monday, true},
		{"monday-thursday matches wednesday", "Monday - Thursday", time.Wednesday, true},
		{"monday-friday matches friday", "Monday - Friday", time.Friday, true},
		{"monday-friday matches tuesday", "Monday - Friday", time.Tuesday, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := scheduleMatchesDay(tt.scheduleName, tt.day)
			if result != tt.want {
				t.Errorf("scheduleMatchesDay(%q, %v) = %v, want %v",
					tt.scheduleName, tt.day, result, tt.want)
			}
		})
	}
}

func TestFormatStopID(t *testing.T) {
	tests := []struct {
		name   string
		stopID string
		want   string
	}{
		{"known station WRP", "940GZZLUWRP", "West Ruislip"},
		{"known station EBY", "940GZZLUEBY", "Ealing Broadway"},
		{"unknown tube station", "940GZZLUXYZ", "XYZ"},
		{"unknown rail station", "910GABC", "ABC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatStopID(tt.stopID)
			if result != tt.want {
				t.Errorf("formatStopID(%q) = %q, want %q", tt.stopID, result, tt.want)
			}
		})
	}
}
