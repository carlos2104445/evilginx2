#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Running Evilginx2 tests...${NC}"

echo -e "${YELLOW}Running Go tests...${NC}"
go test -v -race -coverprofile=coverage.out ./...

if [ -f coverage.out ]; then
    echo -e "${YELLOW}Generating coverage report...${NC}"
    go tool cover -html=coverage.out -o coverage.html
    go tool cover -func=coverage.out
fi

if [ -d "web" ] && [ -f "web/package.json" ]; then
    echo -e "${YELLOW}Running frontend tests...${NC}"
    cd web
    if npm list --depth=0 | grep -q "@testing-library"; then
        npm test -- --coverage --watchAll=false
    else
        echo -e "${YELLOW}No frontend tests configured${NC}"
    fi
    cd ..
fi

echo -e "${GREEN}All tests completed!${NC}"
