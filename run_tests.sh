#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

PORT=5173
SERVER_PID=""

# cleanup
cleanup() {
    if [[ -n "$SERVER_PID" ]]; then
        echo ""
        echo " ===> Stopping backend (PID $SERVER_PID)"
        kill "$SERVER_PID" 2>/dev/null || true
        wait "$SERVER_PID" 2>/dev/null || true
    fi
}
trap cleanup EXIT

# Step 1: build
echo " ===> Building (SQLite)"
make clean build_sqlite

# Step 2: start backend
echo ""
echo " ===> Starting backend on port $PORT"
PORT=$PORT ./main &
SERVER_PID=$!

# Step 3: Wait for the server to be ready
echo " ===> Waiting for server"

status=$(curl -s -o /dev/null -w '%{http_code}' --max-time 5 "http://localhost:$PORT/api/v1/users/login")
if [[ "$status" != 401 ]]; then
    echo " [-!!-] Unexpected response: $status"
    exit 1
fi
echo " ===> Server ready"

# Step 4: API / curl tests
echo ""
echo " ===> Running API tests (make runtests_full)"
make runtests_full
API_STATUS=$?

if [[ $API_STATUS -ne 0 ]]; then
    echo ""
    echo "ERROR: make runtests_full failed (exit $API_STATUS)"
    exit $API_STATUS
fi

echo ""
echo " ===> API tests passed"

# Step 5: Playwright UI tests
echo ""
echo " ===> Running Playwright UI tests"
python -m pytest tests/ -v --tb=short
PLAYWRIGHT_STATUS=$?

# Step 6: Playwright summary
echo ""
if [[ $PLAYWRIGHT_STATUS -eq 0 ]]; then
    echo "All playwright tests passed!"
else
    echo "Playwright tests failed (exit $PLAYWRIGHT_STATUS)"
fi

exit 0
