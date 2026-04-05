#!/bin/bash

echo "Building backend..."
cd backend
go build -o ../bin/server cmd/server/main.go

echo "Building frontend..."
cd ../front
npm run build

echo "Build complete."