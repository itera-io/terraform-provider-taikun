#!/bin/bash

# Initialize attempt counter
attempt=1

# While loop will run until the command succeeds or we've made 3 attempts
while true; do
    # Run the command
    go test . ./taikun -v -run ${ACCEPTANCE_TESTS} -timeout 30m -p 1

    # If the command was successful, break out of the loop
    if [ $? -eq 0 ]; then
        echo "Command succeeded."
        break
    fi

    # If the command failed and we've made 3 attempts, exit with failure status
    if [ $attempt -ge 3 ]; then
        echo "Command failed after 3 attempts."
        exit 1
    fi

    # If the command failed and we've made less than 3 attempts, increment the attempt counter and retry
    echo "Command failed. Attempt #$attempt. Retrying..."
    let "attempt++"
done
