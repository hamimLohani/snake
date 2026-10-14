BINARY_NAME=snake
BIN_DIR=bin
VERSION?=1.0.0
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

.PHONY: all build test lint clean install run man help

all: build

build:
	@mkdir -p $(BIN_DIR)
	@echo "==> Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/snake
	@echo "Build complete: $(BIN_DIR)/$(BINARY_NAME)"

test:
	@echo "==> Running test suite..."
	go test -v -race -cover ./...

lint:
	@echo "==> Running static analysis..."
	go vet ./...

man:
	@echo "==> Updating embedded man page..."
	PAGER=cat man ./man/snake.6 | col -b > cmd/snake/man_embed.txt

install: build
	@echo "==> Installing $(BINARY_NAME) to \$$GOPATH/bin..."
	go install $(LDFLAGS) ./cmd/snake
	@if [ -d /usr/local/share/man/man6 ]; then \
		echo "==> Installing man page to /usr/local/share/man/man6/snake.6..."; \
		install -m 0644 man/snake.6 /usr/local/share/man/man6/snake.6 2>/dev/null || true; \
	fi

clean:
	@echo "==> Cleaning build artifacts..."
	rm -rf $(BIN_DIR)

run: build
	./$(BIN_DIR)/$(BINARY_NAME) play

help:
	@echo "Available Makefile targets:"
	@echo "  build    - Compile the snake binary to bin/snake"
	@echo "  test     - Run all unit and integration tests with race detector"
	@echo "  lint     - Run go vet static checks"
	@echo "  man      - Re-render embedded man page text from troff source"
	@echo "  install  - Install snake binary and man page"
	@echo "  clean    - Remove build artifacts"
	@echo "  run      - Build and run classic mode"
