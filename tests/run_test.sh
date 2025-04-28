#!/bin/bash
set -e

echo "Building project..."
go build ./...

echo "Running tests with verbose output..."
go test -v ./tests

echo "✅ Tests completed successfully."
