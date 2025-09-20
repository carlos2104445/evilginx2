#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Building Evilginx2...${NC}"

if ! command -v go &> /dev/null; then
    echo -e "${RED}Go is not installed. Please install Go 1.21 or later.${NC}"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if ! printf '%s\n%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V -C; then
    echo -e "${RED}Go version $GO_VERSION is too old. Please install Go $REQUIRED_VERSION or later.${NC}"
    exit 1
fi

echo -e "${YELLOW}Building Go backend...${NC}"
go mod download
go build -v -ldflags="-s -w" -o evilginx2 .

if command -v npm &> /dev/null; then
    echo -e "${YELLOW}Building React frontend...${NC}"
    cd web
    npm ci
    npm run build
    cd ..
    echo -e "${GREEN}Frontend built successfully${NC}"
else
    echo -e "${YELLOW}Node.js not found. Skipping frontend build.${NC}"
fi

echo -e "${GREEN}Build completed successfully!${NC}"
echo -e "${GREEN}Binary: ./evilginx2${NC}"
if [ -d "web/dist" ]; then
    echo -e "${GREEN}Frontend: ./web/dist${NC}"
fi
