#!/bin/bash

# Start production backend binary with production config
APP_ENV=production ./bin/server &
BACKEND_PID=$!

# Serve production frontend static files
npx serve dist/front/browser -s -l 4200 &
FRONTEND_PID=$!

echo "Production backend started  (PID: $BACKEND_PID) — config: config.production.yaml"
echo "Production frontend started (PID: $FRONTEND_PID) — http://localhost:4200"
echo "Press Ctrl+C to stop both"

wait $BACKEND_PID $FRONTEND_PID
