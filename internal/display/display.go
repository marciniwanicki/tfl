package display

import (
	"fmt"
	"strings"

	"tfl/internal/tfl"
)

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	red     = "\033[31m"
	green   = "\033[32m"
	yellow  = "\033[33m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
	white   = "\033[37m"
	gray    = "\033[90m"
)

var lineColors = map[string]string{
	"bakerloo":          "\033[48;2;178;99;0m\033[97m",
	"central":           "\033[48;2;220;36;31m\033[97m",
	"circle":            "\033[48;2;255;211;0m\033[30m",
	"district":          "\033[48;2;0;125;50m\033[97m",
	"hammersmith-city":  "\033[48;2;244;169;190m\033[30m",
	"jubilee":           "\033[48;2;161;165;167m\033[30m",
	"metropolitan":      "\033[48;2;155;0;88m\033[97m",
	"northern":          "\033[48;2;0;0;0m\033[97m",
	"piccadilly":        "\033[48;2;0;54;136m\033[97m",
	"victoria":          "\033[48;2;0;160;226m\033[97m",
	"waterloo-city":     "\033[48;2;147;206;186m\033[30m",
	"elizabeth":         "\033[48;2;107;63;160m\033[97m",
	"dlr":               "\033[48;2;0;175;173m\033[97m",
	"london-overground": "\033[48;2;239;123;16m\033[30m",
}

func getLineColor(lineID string) string {
	if color, ok := lineColors[lineID]; ok {
		return color
	}
	return "\033[48;2;100;100;100m\033[97m"
}

func statusColor(severity int) string {
	switch {
	case severity == 10:
		return green
	case severity >= 6 && severity <= 9:
		return yellow
	default:
		return red
	}
}

func PrintLineStatuses(statuses []tfl.LineStatus) {
	fmt.Println()
	fmt.Printf("%s%s TfL Tube Status %s\n\n", bold, white, reset)

	for _, line := range statuses {
		lineCol := getLineColor(line.ID)
		status := line.LineStatuses[0]
		statCol := statusColor(status.StatusSeverity)

		lineName := formatLineName(line.Name)
		fmt.Printf("%s%s%s %s%-20s%s\n",
			lineCol, lineName, reset,
			statCol, status.StatusSeverityDescription, reset)

		if status.Reason != "" {
			reason := wrapText(status.Reason, 60)
			for _, l := range reason {
				fmt.Printf("  %s%s%s\n", gray, l, reset)
			}
		}
	}
	fmt.Println()
}

func PrintDisruptions(disruptions []tfl.Disruption) {
	fmt.Println()
	if len(disruptions) == 0 {
		fmt.Printf("%s%s No current disruptions %s\n\n", bold, green, reset)
		return
	}

	fmt.Printf("%s%s Service Disruptions (%d) %s\n\n", bold, white, len(disruptions), reset)

	for _, d := range disruptions {
		var icon string
		var color string
		switch d.Category {
		case "RealTime":
			icon = "!"
			color = red
		case "PlannedWork":
			icon = "W"
			color = yellow
		default:
			icon = "i"
			color = cyan
		}

		fmt.Printf("%s[%s]%s %s%s%s\n", color, icon, reset, bold, d.CategoryDescription, reset)
		lines := wrapText(d.Description, 70)
		for _, l := range lines {
			fmt.Printf("    %s\n", l)
		}
		fmt.Println()
	}
}

func PrintStopPoints(stops []tfl.StopPoint) {
	if len(stops) == 0 {
		fmt.Printf("%sNo stations found%s\n", yellow, reset)
		return
	}

	fmt.Println()
	fmt.Printf("%s%s Stations found: %s\n\n", bold, white, reset)

	for _, stop := range stops {
		modes := strings.Join(stop.Modes, ", ")
		zone := stop.Zone
		if zone == "" {
			zone = "-"
		}
		fmt.Printf("  %s%-40s%s Zone: %s  [%s]\n", cyan, stop.Name, reset, zone, modes)
		fmt.Printf("  %sID: %s%s\n\n", gray, stop.ID, reset)
	}
}

const lineNameWidth = 14

func formatLineName(name string) string {
	if len(name) > lineNameWidth {
		return " " + name[:lineNameWidth-2] + ".. "
	}
	// Pad to fixed width
	padded := " " + name
	for len(padded) < lineNameWidth+1 {
		padded += " "
	}
	return padded + " "
}

func PrintArrivals(arrivals []tfl.Arrival, stationName string) {
	fmt.Println()
	if len(arrivals) == 0 {
		fmt.Printf("%sNo arrivals found for %s%s\n\n", yellow, stationName, reset)
		return
	}

	fmt.Printf("%s%s Departures from %s %s\n\n", bold, white, stationName, reset)

	for _, arr := range arrivals {
		lineCol := getLineColor(arr.LineID)
		mins := arr.TimeToStation / 60
		departureTime := arr.ExpectedArrival.Local().Format("15:04")

		var timeStr string
		switch {
		case mins == 0:
			timeStr = fmt.Sprintf("%sDue%s", green+bold, reset)
		case mins == 1:
			timeStr = fmt.Sprintf("%s1 min%s", green, reset)
		case mins < 60:
			timeStr = fmt.Sprintf("%d mins", mins)
		default:
			hours := mins / 60
			remainMins := mins % 60
			if remainMins == 0 {
				timeStr = fmt.Sprintf("%dh", hours)
			} else {
				timeStr = fmt.Sprintf("%dh %dm", hours, remainMins)
			}
		}

		platform := arr.PlatformName
		if platform == "" {
			platform = "-"
		}

		lineName := formatLineName(arr.LineName)
		fmt.Printf("%s%s%s  %s%s%s  %-8s  %s%-28s%s  %s%s%s\n",
			lineCol, lineName, reset,
			cyan, departureTime, reset,
			timeStr,
			bold, arr.DestinationName, reset,
			gray, platform, reset)
	}
	fmt.Println()
}

func wrapText(text string, width int) []string {
	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return lines
	}

	currentLine := words[0]
	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)
	return lines
}
