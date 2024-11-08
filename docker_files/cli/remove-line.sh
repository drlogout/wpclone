#!/bin/bash

# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 'line of text' and filename"
    exit 1
fi

line=$1
file=$2

if [ ! -f "$file" ]; then
    # Create the file if it does not exist
    touch "$file"
fi

# Remove the line if it exists
if grep -Fxq "$line" "$file"; then
    sed -i "/$line/d" "$file"
fi

# Remove the anotation if it exists
if grep -Fxq "# Added by wpclone" "$file"; then
    sed -i "/# Added by wpclone/d" "$file"
fi
