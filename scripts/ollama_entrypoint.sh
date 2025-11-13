#!/bin/bash

# Start the Ollama server in the background
/bin/ollama serve &

# Record the Process ID (PID) of the background process
pid=$!

# Wait for the Ollama server to become responsive
echo "Waiting for Ollama server to start..."
until ollama list > /dev/null 2>&1; do
  sleep 1
done
echo "Ollama server is running."

# Pull the desired model using the value from an environment variable
MODEL_NAME=${OLLAMA_MODEL:-llama3} # Default to 'llama3' if the variable is not set

echo "Pulling model: $MODEL_NAME"
ollama pull "$MODEL_NAME"

echo "Model pull complete."

# Wait for the main Ollama server process to finish (which it won't, it runs forever)
# This keeps the container running with the server active.
wait $pid