# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

TfL CLI - A Go command-line tool for accessing Transport for London services including real-time departures, tube line status, and service disruptions. Uses the public TfL API (api.tfl.gov.uk).

## Build Commands

```bash
make build      # Compile binary (go build -o tfl .)
make test       # Run tests (go test -v ./...)
make lint       # Run linter (golangci-lint run)
make format     # Format code (go fmt ./...)
make clean      # Remove compiled binary
```

## Architecture

The codebase follows a clean separation between CLI commands, API client, and display logic:

- **cmd/**: Cobra command implementations
  - `root.go`: Root command, global flags, client initialization via PersistentPreRun
  - `departures.go`: Station departures and search commands
  - `status.go`: Tube line status command
  - `disruptions.go`: Service disruptions command

- **internal/tfl/**: API client layer
  - `client.go`: HTTP client with methods for TfL API endpoints
  - `types.go`: Data models (LineStatus, StopPoint, Arrival, etc.)

- **internal/display/**: Terminal output formatting with ANSI colors matching actual TfL line branding

## Key Patterns

- API key passed via `--key` flag or `TFL_APP_KEY` environment variable
- Client created once in root command's PersistentPreRun and shared across subcommands
- All API methods return typed results with errors propagated up the call stack
- Arrivals are auto-sorted by time to station

## Running the CLI

```bash
./tfl status                              # All tube line statuses
./tfl disruptions                         # Current service disruptions
./tfl departures "Liverpool Street"       # Departures from station
./tfl departures paddington central -n 5  # Filter by line, limit results
./tfl search "King's Cross"               # Search for stations
```
