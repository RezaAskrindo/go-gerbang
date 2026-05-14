#!/bin/bash

# Kill any existing apigateway-9000 processes
pkill -f apigateway-9000

# Wait briefly to ensure the process has time to terminate
sleep 2

# Check if the binary exists and is executable
if [ ! -x PATH/apigateway-9000 ]; then
  echo "Error: binary 'PATH/apigateway-9000' not found or not executable."
  exit 1
fi

# Start the binary in the background and disown it
PATH/apigateway-9000 & disown

echo "apigateway-9000 restarted successfully."