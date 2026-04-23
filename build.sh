#!/bin/bash
set -e
cd "$(dirname "$0")"

echo "Building futuapi4go-demo..."

go build ./...

echo "Build successful."
