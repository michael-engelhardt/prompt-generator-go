#!/bin/bash

#!/bin/bash

# Set the Go binary name
working_dir=$(pwd)
binary_path="$working_dir/bin"
binary_name="prompt-generator-linux-amd64"
full_binary_path="${binary_path}/${binary_name}"
global_binary_name="prompt-generator"


# Build for Linux
GOOS=linux GOARCH=amd64 go build -o ${full_binary_path}

# Move the binary to the local bin folder
mv ${full_binary_path} /usr/local/bin/${global_binary_name}


# After installation run with:
# prompt-generator