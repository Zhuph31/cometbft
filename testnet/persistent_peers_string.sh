#!/bin/bash

# Check if the required argument is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <ip1> <ip2> <ip3> ..."
    exit 1
fi

# Command to run on each IP
BASE_COMMAND="cometbft show_node_id --home ./mytestnet/node"

# Initialize an array to store results
PERSISTENT_PEERS=""

port=26656

# Iterate through provided IPs
for i in "$@"; do
    IP="$i"
    NODE_IDX=$((i - 1)) # Adjust for zero-based indexing

    echo "Getting ID of $IP (node $NODE_IDX)..."

    # Run the command on the current IP and capture the result
    ID=$($BASE_COMMAND$NODE_IDX)

    # Store the result in the array
    PERSISTENT_PEERS+="$ID@$IP:$port"
    ((port++))

    # Add a comma if not the last IP
    if [ $i -lt $# ]; then
        PERSISTENT_PEERS+=","
    fi
done

echo "$PERSISTENT_PEERS"
