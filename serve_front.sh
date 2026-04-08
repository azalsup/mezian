#!/bin/bash

echo "Copying categories.json to frontend public/..."
cp backend/data/categories.json front/public/categories.json

echo "Starting frontend in development mode..."
cd front
ng serve