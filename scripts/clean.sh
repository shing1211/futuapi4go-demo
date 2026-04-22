#!/bin/bash
echo "Cleaning build artifacts..."

rm -f cmd/demo/futuapi4go-demo
rm -f go.work
rm -f go.work.sum

echo "Done."
