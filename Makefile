# Project metadata
BINARY_NAME := pihole-mcp
DIST_DIR := dist

# Version info (optional but useful later)
VERSION ?= $(shell git describe --tags --dirty --always 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -X main.version=$(VERSION) \
           -X main.commit=$(COMMIT) \
           -X main.date=$(DATE)

GO := go
GOFLAGS :=

# Default target
.PHONY: all
all: build

# --------------------
# Local development
# --------------------

.PHONY: build
build:
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)

.PHONY: run
run:
	$(GO) run .

.PHONY: test
test:
	$(GO) test ./...

.PHONY: fmt
fmt:
	$(GO) fmt ./...

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: lint
lint:
	@command -v golangci-lint >/dev/null || (echo "golangci-lint not installed"; exit 1)
	golangci-lint run

# --------------------
# Cross-compilation
# --------------------

.PHONY: dist
dist: clean-dist dist-linux-amd64 dist-linux-arm64 dist-linux-armv7

.PHONY: dist-linux-amd64
dist-linux-amd64:
	mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
	-o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64

.PHONY: dist-linux-arm64
dist-linux-arm64:
	mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 \
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
	-o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64

.PHONY: dist-linux-armv7
dist-linux-armv7:
	mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 \
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" \
	-o $(DIST_DIR)/$(BINARY_NAME)-linux-armv7

# --------------------
# Housekeeping
# --------------------

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

.PHONY: clean-dist
clean-dist:
	rm -rf $(DIST_DIR)

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: deps
deps:
	$(GO) mod download
