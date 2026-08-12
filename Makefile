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
	rm -f [0-9][0-9]_
	rm -f [0-9][0-9][0-9]_
	rm -f [0-9][0-9][0-9][0-9]_*

examples:
	@ls -1 examples/ | sort -n

run:
	@if [ -z "$(EXP)" ]; then \
		echo "Usage: make run EXP=00"; \
		echo "Available: 00-110"; \
	else \
		go run ./examples/$(EXP); \
	fi

run-all:
	@total=0; failed=0; failed_list=""; \
	for dir in examples/*/; do \
		name=$$(basename $$dir); \
		case "$$name" in pkg|graphify-out) continue;; esac; \
		total=$$((total + 1)); \
		if go build -o /dev/null ./$$dir 2>/dev/null; then \
			printf "  [PASS] %s\n" "$$name"; \
		else \
			failed=$$((failed + 1)); \
			failed_list="$$failed_list $$name"; \
			printf "  [FAIL] %s\n" "$$name"; \
		fi; \
	done; \
	echo "-----------------------------------------"; \
	echo "  Built $$total examples, $$failed failed."; \
	if [ -n "$$failed_list" ]; then \
		echo "  Failed:$$failed_list"; \
		exit 1; \
	fi

deps:
	go mod download
	go mod tidy