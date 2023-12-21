#!/bin/bash

# Set the Go binary name
working_dir=$(pwd)
binary_path="$working_dir/bin"
binary_name="prompt-generator.exe"

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o ${binary_path}/${binary_name}
