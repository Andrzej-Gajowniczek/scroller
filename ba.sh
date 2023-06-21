#!/bin/bash

# Extract all go get comments from Go source file
get_comments=$(grep -Eo '// go get [^[:space:]]+' 3d.go)

# Loop through each go get comment and execute the go get command
while IFS= read -r line; do
    # Extract the package name from the go get comment
    package=$(echo "$line" | awk '{print $3}')
    
    # Execute go get command to install the package
    go get "$package"
done <<< "$get_comments"

