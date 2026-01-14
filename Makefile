# Makefile for termidoro

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
BINARY_NAME=termidoro
MAN_PAGE=termidoro.1

# Default target
all: build

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_NAME)

# Run tests
test:
	$(GOTEST) ./...

# Install the binary and man page
install: build
	install -m 755 $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	install -m 644 $(MAN_PAGE) /usr/local/share/man/man1/$(MAN_PAGE).1

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

.PHONY: all build test install clean
