#!/bin/bash
set -e
cd "$(dirname "$0")/.."

echo "Building futuapi4go-demo..."

OUTPUT=cmd/demo/futuapi4go-demo
go build -o "$OUTPUT" ./cmd/demo

echo "Build successful: $OUTPUT"
