.PHONY: build clean test install release-dry release

# Variables
BINARY_NAME=zellijinator
MAIN_PACKAGE=.
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-s -w -X main.version=${VERSION}"

# Build for current platform
build:
	go build ${LDFLAGS} -o ${BINARY_NAME} ${MAIN_PACKAGE}

# Build for all platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-amd64 ${MAIN_PACKAGE}
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-arm64 ${MAIN_PACKAGE}
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-amd64 ${MAIN_PACKAGE}

# Run tests
test:
	go test -v ./...

# Install locally
install: build
	sudo mv ${BINARY_NAME} /usr/local/bin/

# Clean build artifacts
clean:
	rm -f ${BINARY_NAME}
	rm -rf dist/

# Run goreleaser in snapshot mode (doesn't publish)
release-dry:
	goreleaser release --snapshot --skip=publish --clean

# Create a new release (requires a tag)
release:
	goreleaser release --clean

# Format code
fmt:
	go fmt ./...

# Run linters
lint:
	golangci-lint run

# Tidy dependencies
tidy:
	go mod tidy