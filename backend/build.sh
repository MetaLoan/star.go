#!/bin/bash

# Star Backend Build Script
# Swiss Ephemeris is the ONLY data source for all astronomical calculations

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

echo "🌟 Building Star API with Swiss Ephemeris..."

# Check if Swiss Ephemeris library exists
if [ ! -f "lib/libswe.a" ]; then
    echo "❌ Swiss Ephemeris library not found at lib/libswe.a"
    echo "   Please install Swiss Ephemeris first."
    exit 1
fi

if [ ! -f "include/swephexp.h" ]; then
    echo "❌ Swiss Ephemeris headers not found at include/"
    echo "   Please install Swiss Ephemeris first."
    exit 1
fi

# Set CGO environment variables
export CGO_ENABLED=1
export CGO_LDFLAGS="-L${SCRIPT_DIR}/lib"
export CGO_CFLAGS="-I${SCRIPT_DIR}/include"

# Build with swe tag
echo "🔧 Compiling with CGO and Swiss Ephemeris..."
go build -tags swe -o bin/star-api .

echo "✅ Build successful!"
echo "📦 Binary: bin/star-api"
echo ""
echo "🚀 Run with: ./bin/star-api"
echo ""
echo "⚠️  Important: Swiss Ephemeris is the ONLY data source for astronomical calculations."

