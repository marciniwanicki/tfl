# TfL CLI

A command-line tool for accessing Transport for London services including real-time departures, tube line status, and service disruptions.

## Features

- **Real-time departures** from any tube, Elizabeth line, DLR, or Overground station
- **Tube & Elizabeth line status** with service alerts and disruption details
- **Fuzzy filtering** by line, destination, or platform
- **Time-based filtering** for departures at specific times
- **Timetable support** for tube lines (scheduled departures hours ahead)
- **Colour-coded output** matching official TfL line branding

## Installation

```bash
# Clone and build
git clone https://github.com/yourusername/tfl.git
cd tfl
make build

# Optional: move to PATH
mv tfl /usr/local/bin/
```

## Usage

### Line Status

```bash
# Show status for all tube and Elizabeth lines
tfl status
```

### Departures

```bash
# Show departures from a station
tfl departures "liverpool street"
tfl departures paddington

# Limit number of results
tfl departures "kings cross" -n 5

# Filter by line
tfl departures paddington central
tfl departures stratford elizabeth

# Filter by time (shows departures at or after specified time)
tfl departures paddington -t 14:30

# Fuzzy match on line, destination, or platform
tfl departures "kings cross" -m "eastbound"
tfl departures stratford -m "heathrow"
tfl departures paddington -m "ealing central"

# Combine filters
tfl departures "liverpool street" -m "westbound" -t 18:00 -n 10
```

### Search Stations

```bash
# Find stations by name
tfl search "victoria"
tfl search "kings cross"
```

### Disruptions

```bash
# Show current service disruptions
tfl disruptions
```

## API Key

The TfL API works without a key for basic usage, but you may want to register for higher rate limits:

```bash
# Set via environment variable
export TFL_APP_KEY=your_api_key

# Or pass via flag
tfl status --key your_api_key
```

Register for a free API key at [TfL API Portal](https://api-portal.tfl.gov.uk/).

## Examples

### Morning commute check
```bash
# Check line status and next few trains
tfl status
tfl departures "finsbury park" -n 5
```

### Planning ahead
```bash
# Tube lines support timetable lookups for future times
tfl departures "bank" central -t 18:30

# Elizabeth line only shows real-time (~30 min ahead)
tfl departures stratford elizabeth -t 11:00
```

### Filtering busy stations
```bash
# Show only eastbound Piccadilly line trains
tfl departures "kings cross" -m "piccadilly eastbound"

# Show trains to a specific destination
tfl departures paddington -m "heathrow terminal"
```

## Build Commands

```bash
make build      # Compile binary
make test       # Run tests
make lint       # Run linter
make format     # Format code
make clean      # Remove binary
```

## Limitations

- **Elizabeth line timetables**: The TfL API doesn't provide timetable data for Elizabeth line. Only real-time arrivals (~30 minutes ahead) are available.
- **National Rail**: Departures from National Rail services at shared stations may appear but are not fully supported.

## License

TfL CLI is released under version 2.0 of the [Apache License](LICENSE).
