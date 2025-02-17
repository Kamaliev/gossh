#!/bin/bash

OUTPUT_DIR="build"
APP=gossh
mkdir -p $OUTPUT_DIR

echo "🔨 Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o $OUTPUT_DIR/$APP-linux-amd64 .

echo "🔨 Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o $OUTPUT_DIR/$APP-mac-amd64 .

echo "🔨 Building for macOS (arm64)..."
GOOS=darwin GOARCH=arm64 go build -o $OUTPUT_DIR/$APP-mac-arm64 .

echo "✅ Build completed. Binaries are in the '$OUTPUT_DIR' directory."
