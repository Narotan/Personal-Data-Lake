#!/bin/bash

set -e

echo "Building aw-client..."
cd "$(dirname "$0")/.."
go build -o bin/aw-client ./cmd/aw-client

echo "Build complete: bin/aw-client"
echo "To run manually: ./bin/aw-client -minutes 5"

