#!/bin/bash
set -e
cd "$(dirname "$0")"

echo "Cleaning build artifacts..."

rm -f futuapi4go-demo
rm -f examples/futuapi4go-demo

echo "Done."
