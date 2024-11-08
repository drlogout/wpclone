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

if ! grep -Fxq "$line" "$file"; then
    { echo ""; echo ""; echo "# Added by wpclone"; } >>"$file"
    echo "$line" >>"$file"
fi
