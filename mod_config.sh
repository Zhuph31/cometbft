#!/bin/bash

# Function to find line number containing search string and replace text
replace_text() {
    local file="$1"               # File to be modified
    local search_string="$2"      # String to search for
    local old_string="$3"         # String to be replaced
    local replacement_string="$4" # Text to replace with

    # Find the line number containing the search string
    local line_number=$(grep -n "$search_string" "$file" | cut -d':' -f1)

    if [ -n "$line_number" ]; then
        # Store the old line
        local old_line=$(sed -n "${line_number}p" "$file")

        # Replace text on the found line
        sed -i "${line_number}s/$old_string/$replacement_string/" "$file"

        # Store the new line
        local new_line=$(sed -n "${line_number}p" "$file")

        # Print old and new lines along with file path
        echo "File: $file"
        echo "Old line: $old_line"
        echo "New line: $new_line"
        echo "Replacement done on line $line_number!"
    else
        echo "Search string not found in the file."
    fi
}

# modify the config for all the docker nodes, so that cometbft would work as we expect them to
cd build

# Iterate through each folder in the current directory
for folder in */; do
    # Remove the trailing slash from the folder name
    folder=${folder%/}

    # Perform your action here, for example, print the folder name
    echo "Processing folder: $folder"

    cd $folder/config

    replace_text "config.toml" "log_level" "info" "debug"
    pwd

    cd ../..

    # Add your action here
    # Example: Perform some operation on the folder
    # Example: Call another script or command with the folder as argument
done
