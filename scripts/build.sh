#!/bin/bash
echo "Building futuapi4go-demo..."

OUTPUT=cmd/demo/futuapi4go-demo
go build -o "$OUTPUT" ./cmd/demo

if [ $? -eq 0 ]; then
    echo "Build successful: $OUTPUT"
else
    echo "Build failed!"
    exit 1
fi
