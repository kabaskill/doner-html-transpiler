#!/bin/bash

# Deployment script for D.Ö.N.E.R
echo "🥙 Deploying D.Ö.N.E.R to Render..."

# Build frontend
echo "📦 Building frontend..."
cd frontend
npm install
npm run build
cd ..

# Test backend compilation
echo "🔨 Testing backend build..."
cd backend
go mod tidy
go build -o main .
rm -f main  # Remove the test binary
cd ..

echo "✅ Build test successful!"
echo "🚀 Push to your GitHub repository and connect it to Render."
echo ""
echo "Render deployment steps:"
echo "1. Go to https://render.com"
echo "2. Connect your GitHub repository"
echo "3. Choose 'Web Service'"
echo "4. Select Docker environment"
echo "5. Use the included Dockerfile"
echo "6. Deploy!"
