#!/bin/bash
echo "Upgrading dependencies..."

go get -u github.com/shing1211/futuapi4go@latest
go mod tidy

if [ $? -eq 0 ]; then
    echo "Upgrade successful!"
else
    echo "Upgrade failed!"
    exit 1
fi
