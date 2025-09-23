# Makefile for aws-s3-uploader

# Default target
all: build

# Build for current platform
build:
	go build -o aws-s3-uploader .

# Cross-compile for common platforms
cross-compile:
	mkdir -p dist/darwin-amd64 dist/darwin-arm64 dist/linux-amd64 dist/linux-arm64 dist/windows-amd64
	GOOS=darwin GOARCH=amd64 go build -o dist/darwin-amd64/aws-s3-uploader .
	GOOS=darwin GOARCH=arm64 go build -o dist/darwin-arm64/aws-s3-uploader .
	GOOS=linux GOARCH=amd64 go build -o dist/linux-amd64/aws-s3-uploader .
	GOOS=linux GOARCH=arm64 go build -o dist/linux-arm64/aws-s3-uploader .
	GOOS=windows GOARCH=amd64 go build -o dist/windows-amd64/aws-s3-uploader.exe .

# Clean build artifacts
clean:
	rm -rf dist/
	rm -f aws-s3-uploader
	rm -f aws-s3-uploader.exe

# Install dependencies
deps:
	go mod tidy

.PHONY: all build cross-compile clean deps