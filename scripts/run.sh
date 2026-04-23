#!/bin/bash
set -e
cd "$(dirname "$0")/.."

ADDR="${FUTU_ADDR:-127.0.0.1:11111}"

echo "Running futuapi4go-demo (OpenD: $ADDR)..."
echo

go run ./cmd/demo/main.go
