#!/bin/bash
set -e
cd "$(dirname "$0")"

echo "Upgrading dependencies..."

go get -u github.com/shing1211/futuapi4go@latest
go mod tidy

echo "Upgrade successful!"
