#!/bin/bash

# Start production backend binary
./bin/server &
BACKEND_PID=$!

# Serve production frontend static files
npx serve dist/front/browser -s -l 4200 &
FRONTEND_PID=$!

echo "Production backend started (PID: $BACKEND_PID)"
echo "Production frontend started (PID: $FRONTEND_PID)"
echo "Press Ctrl+C to stop both"

# Wait for both
wait $BACKEND_PID $FRONTEND_PID