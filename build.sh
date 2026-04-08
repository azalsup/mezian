#!/bin/bash
set -e

PROD_API_URL="https://api.daba.incipyo.com/api/v1"
DEV_API_URL="/api/v1"
APP_CONFIG="front/public/app-config.json"

# ── Frontend ────────────────────────────────────────────────────────────────

echo "Writing production app-config.json..."
printf '{\n  "apiBaseUrl": "%s"\n}\n' "$PROD_API_URL" > "$APP_CONFIG"

echo "Copying categories.json to frontend public/..."
cp backend/data/categories.json front/public/categories.json

echo "Building frontend (SSG)..."
cd front
npx ng build --base-href=/ --output-mode static
cd ..

echo "Restoring dev app-config.json..."
printf '{\n  "apiBaseUrl": "%s"\n}\n' "$DEV_API_URL" > "$APP_CONFIG"

# ── Backend ─────────────────────────────────────────────────────────────────

echo "Building backend..."
mkdir -p bin
cd backend
go build -o ../bin/server cmd/server/main.go
cd ..

# ── Done ────────────────────────────────────────────────────────────────────

echo ""
echo "Build complete."
echo "  Frontend → front/dist/front/browser/  (static HTML + assets)"
echo "  Backend  → bin/server"
echo ""
echo "To start production: APP_ENV=production ./start.sh"
echo "  (override secrets with: JWT_SECRET=... DB_PATH=... APP_ENV=production ./start.sh)"
