#!/bin/zsh

# Build
go build -o media-analyzer cmd/analyzer/main.go

# Define the media path
MEDIA_PATH="/Volumes/Media Drive/Photos/Family"

# Run without web interface
./media-analyzer --path "${MEDIA_PATH}"

# Run with web interface
#./media-analyzer --path "${MEDIA_PATH}" --web