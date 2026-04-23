#!/bin/bash
set -e
cd "$(dirname "$0")"

ADDR="${FUTU_ADDR:-127.0.0.1:11111}"
EXAMPLE="${1:-00_connect}"

echo "Running futuapi4go-demo/examples/$EXAMPLE (OpenD: $ADDR)..."
echo

go run ./examples/"$EXAMPLE"
