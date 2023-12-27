#!/usr/bin/env sh

set -e

# Create temporary directory
tmpfile=$(mktemp)
# Build the project
(cd $(dirname "$0") && go build -o "$tmpfile" ./cmd/dnsd )
# Run the project with config
exec "$tmpfile" "-config=$(dirname "$0")/config.json"