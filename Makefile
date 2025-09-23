# Makefile for cloudflare-r2-uploader

# Default target
all: build

# Build for current platform
build:
	go build -o cloudflare-r2-uploader .

# Cross-compile for common platforms
cross-compile:
	mkdir -p dist/darwin-amd64 dist/darwin-arm64 dist/linux-amd64 dist/linux-arm64 dist/windows-amd64
	GOOS=darwin GOARCH=amd64 go build -o dist/darwin-amd64/cloudflare-r2-uploader .
	GOOS=darwin GOARCH=arm64 go build -o dist/darwin-arm64/cloudflare-r2-uploader .
	GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/cloudflare-r2-uploader .
	GOOS=linux GOARCH=arm64 go build -o dist/linux-arm64/cloudflare-r2-uploader .
	GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/cloudflare-r2-uploader.exe .

# Clean build artifacts
clean:
	rm -rf dist/
	rm -f cloudflare-r2-uploader
	rm -f cloudflare-r2-uploader.exe

# Install dependencies
deps:
	go mod tidy

.PHONY: all build cross-compile clean deps