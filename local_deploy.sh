#!/bin/bash
# deploy.sh - Simple deployment script for your Go "daba" service

set -euo pipefail  # Exit on error, undefined vars, and pipe failures

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "🚀 Starting deployment from: $SCRIPT_DIR"

# 1. Go to backend directory
cd "$SCRIPT_DIR/backend" || { echo "❌ Cannot cd to backend directory"; exit 1; }

# 2. Build the Go binary
echo "🔨 Building Go server..."
go build -o ../bin/server cmd/server/main.go

# 3. Stop the service
echo "⏹️  Stopping daba service..."
sudo systemctl stop daba

# 4. Copy the new binary
echo "📦 Copying binary to /opt/daba/bin/"
sudo cp ../bin/server /opt/daba/bin/

# 5. Fix permissions
echo "🔑 Setting ownership to dabauser:dabauser..."
sudo chown -R dabauser:dabauser /opt/daba/bin/

# 6. Start the service again
echo "▶️  Starting daba service..."
sudo systemctl start daba

# Optional: Show status
echo "✅ Deployment finished. Checking service status..."
sudo systemctl status daba --no-pager -l
