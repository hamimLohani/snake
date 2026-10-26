BINARY_NAME=snake
BIN_DIR=bin
DIST_DIR=dist
VERSION?=$(shell git describe --tags --always 2>/dev/null || echo "1.0.9")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: all build test lint clean install run man release help

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

release:
	@echo "==> Cross-compiling release binaries for $(VERSION)..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/snake-linux-amd64 ./cmd/snake
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/snake-linux-arm64 ./cmd/snake
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(DIST_DIR)/snake-darwin-arm64 ./cmd/snake
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/snake-darwin-amd64 ./cmd/snake
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(DIST_DIR)/snake-windows-amd64.exe ./cmd/snake
	@echo "==> Packaging archives..."
	cd $(DIST_DIR) && tar -czf snake-$(VERSION)-linux-amd64.tar.gz snake-linux-amd64
	cd $(DIST_DIR) && tar -czf snake-$(VERSION)-linux-arm64.tar.gz snake-linux-arm64
	cd $(DIST_DIR) && tar -czf snake-$(VERSION)-darwin-arm64.tar.gz snake-darwin-arm64
	cd $(DIST_DIR) && tar -czf snake-$(VERSION)-darwin-amd64.tar.gz snake-darwin-amd64
	cd $(DIST_DIR) && shasum -a 256 *.tar.gz > checksums.txt
	@echo "==> Release artifacts ready in $(DIST_DIR)/:"
	@ls -lh $(DIST_DIR)/

clean:
	@echo "==> Cleaning build artifacts..."
	rm -rf $(BIN_DIR) $(DIST_DIR)

run: build
	./$(BIN_DIR)/$(BINARY_NAME) play

help:
	@echo "Available Makefile targets:"
	@echo "  build    - Compile the snake binary to bin/snake"
	@echo "  test     - Run all unit and integration tests with race detector"
	@echo "  lint     - Run go vet static checks"
	@echo "  man      - Re-render embedded man page text from troff source"
	@echo "  install  - Install snake binary and man page"
	@echo "  release  - Cross-compile release binaries and packages in dist/"
	@echo "  clean    - Remove build artifacts"
	@echo "  run      - Build and run classic mode"
