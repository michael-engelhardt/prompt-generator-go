#!/bin/bash

# Set the Go binary name
binary_path="bin"
binary_name="prompt-generator"

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o ${binary_path}/${binary_name}

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o ${binary_path}/${binary_name}.exe