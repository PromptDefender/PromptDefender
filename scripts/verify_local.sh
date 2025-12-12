#!/bin/bash
set -e

# Configuration
PROJECT_ROOT=$(pwd)
SERVER_DIR="$PROJECT_ROOT/cmd/server"
SERVER_BIN="$SERVER_DIR/server"
TEST_DIR="$PROJECT_ROOT/test/integration_test_harness"
PORT=8081

echo "Building server..."
cd "$SERVER_DIR"
go build -o server
cd "$PROJECT_ROOT"

echo "Starting server in background..."
export PORT=$PORT
# Mock AWS region if needed for SDK initialization
export AWS_REGION="eu-west-1" 
export CACHE_TABLE_NAME="test-cache-table"
export USERS_TABLE="test-users-table"
export open_ai_api_key="sk-dummy-key-for-testing-startup"
# Start server and redirect logs
"$SERVER_BIN" > server.log 2>&1 &
SERVER_PID=$!
echo "Server PID: $SERVER_PID"

# Wait for server to start
sleep 3

echo "Running integration tests..."
export URL="http://localhost:$PORT"
export DEFENDER_API_KEY="dummy-key" # Server currently doesn't validate key against a DB in this simple impl, or we assume it just accepts any for now unless middleware checks it specifically.
# Note: The api.gen.go middleware adds "ApiKeyAuth.Scopes" but actual validation logic depends on implementation.
# In `main.go`, `ApiKeyAuth` is not explicitly wired with a validator middleware that rejects requests yet? 
# Wait, `api.gen.go` generates `ServerInterfaceWrapper` which puts scopes in context, but doesn't mandate auth unless middleware is added.
# In `main.go`, we used `HandlerFromMux(server, r)`.
# `HandlerFromMux` in `api.gen.go` calls `HandlerWithOptions`.
# `HandlerWithOptions` iterates middlewares.
# We didn't add an auth middleware in `main.go`. So it should be open or accept anything.
# But `GetUser` logic in server might check `x-api-key`.

cd "$TEST_DIR"
# Run tests. Use -count=1 to avoid caching.
go test -count=1 -v ./...

TEST_EXIT_CODE=$?

echo "Stopping server..."
kill "$SERVER_PID"

if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo "Tests PASSED"
else
    echo "Tests FAILED"
    echo "Server Logs:"
    cat "$PROJECT_ROOT/server.log"
    exit 1
fi
