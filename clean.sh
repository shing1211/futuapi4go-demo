#!/bin/bash
set -e
cd "$(dirname "$0")"

echo "Cleaning build artifacts..."

rm -f cmd/demo/futuapi4go-demo

echo "Done."
