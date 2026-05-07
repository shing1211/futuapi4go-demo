.PHONY: build vet test clean examples help

help:
	@echo "futuapi4go-demo Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  build      - Build all examples"
	@echo "  vet        - Run go vet for linting"
	@echo "  examples   - List all example directories"
	@echo "  run EXP=00 - Run example (e.g., make run EXP=01)"
	@echo "  clean      - Remove compiled binaries"

build:
	go build ./...

vet:
	go vet ./...

test:
	go test ./...

clean:
	rm -f examples/*/*.exe
	rm -f *.exe

examples:
	@ls -1 examples/ | sort -n

run:
	@if [ -z "$(EXP)" ]; then \
		echo "Usage: make run EXP=00"; \
		echo "Available: 00-80"; \
	else \
		go run ./examples/$(EXP); \
	fi

run-all:
	@for dir in examples/*/; do \
		name=$$(basename $$dir); \
		echo "Building $$name..."; \
		go build -o /dev/null ./$$dir || echo "  FAILED"; \
	done
	@echo "Done."

deps:
	go mod download
	go mod tidy