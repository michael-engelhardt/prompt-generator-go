#!/bin/bash

# Set the Go binary name
working_dir=$(pwd)
binary_path="$working_dir/bin"
binary_name="prompt-generator-windows-amd64.exe"
full_binary_path="${binary_path}/${binary_name}"

# Run the binary
${full_binary_path}