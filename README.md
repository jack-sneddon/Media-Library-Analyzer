# Media Library Analyzer

A small tool for analyzing and organizing digital photo libraries based on my media organization methodology. It helps identify gaps and inconsistencies in the photo collection organization.

## Directory Organization

The analyzer expects photos to be organized in a year/month hierarchy:

```bash
Photos/
├── 2024/
│   ├── 01-January/
│   ├── 02-February/
│   └── ...
├── 2023/
│   ├── 01-January/
│   ├── 02-February/
│   └── ...
└── ...
```

Each year directory contains 12 month directories, named using the pattern `MM-Month` (e.g., "01-January", "02-February").

## Features

- **Library Analysis**
  - Scans directory structure
  - Counts files per month
  - Identifies missing months
  - Flags months with low content (< 10 files)
  - Tracks total file counts per year

- **Interactive HTML Report**
  - Visual dashboard of your library
  - Color-coded status indicators
  - Real-time filtering capabilities
  - Year and month breakdowns
  - File count statistics

- **Prioritized Console Summary**
  - High priority issues (last 2 years)
  - Medium priority issues (3-5 years ago)
  - Historical gaps (older than 5 years)
  - Actionable recommendations

## Installation

```bash
# Clone the repository
git clone https://github.com/jack-sneddon/Media-Library-Analyzer.git

# Build the application
cd Media-Library-Analyzer
go build -o media-analyzer cmd/analyzer/main.go
```

## Usage

The analyzer can be run in two modes: Console (default) and Web Interface.

### Console Mode

Generates a text report in the terminal and saves an HTML report file:

```bash
# Basic usage
./media-analyzer --path "/path/to/photos"

# Example with path containing spaces
./media-analyzer --path "/Volumes/Media Drive/Photos"
```

Console mode will:

1. Analyze your photo library
2. Display a prioritized summary in the terminal
3. Save a static HTML report as `library_analysis.html`

### Web Interface Mode

Starts a local web server with an interactive dashboard:

```bash
# Basic usage
./media-analyzer --web --path "/path/to/photos"

# Specify custom port (default is 8080)
./media-analyzer --web --path "/path/to/photos" --port 8888

# Example with actual media library path
./media-analyzer --web --path "/Volumes/Media Drive/Photos/Family"
```

Web mode will:

1. Analyze your photo library
2. Start a local web server
3. Open your default browser to the dashboard
4. Display an interactive interface with:
   - Real-time filtering
   - Color-coded status indicators
   - Detailed statistics
   - Year and month breakdowns
5. Keep running until you press Ctrl+C

### Available Options

```bash
--path    Path to your photo library (required)
--web     Start web interface (optional)
--port    Specify web server port (optional, default: 8080)
```

## Report Types

### Console Summary

The console output provides a prioritized list of issues:

- Recent gaps (high priority)
- Medium-term gaps
- Historical gaps
- Recommendations for addressing issues

### HTML Dashboard

The interactive dashboard allows you to:

- Filter by decade, year, or month
- Focus on specific status types (missing/light/normal)
- Set minimum file count thresholds
- Visualize patterns and trends

## Status Definitions

- **Normal**: 10 or more files in a month
- **Light**: 1-9 files in a month
- **Missing**: No files in a month

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Author

Created by Jack Sneddon

## Acknowledgments

- Inspired by the need for better photo library organization
- Built with Go and modern web technologies

## Requirements

- Go 1.21 or higher
- Modern web browser for HTML report viewing

## Roadmap

Future enhancements may include:

- File type analysis
- Duplicate detection
- Metadata extraction
- Batch organization tools
- More detailed statistics

## Support

For issues and feature requests, please use the GitHub issue tracker.