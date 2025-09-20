#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

IMAGE_NAME="evilginx2"
TAG=${1:-latest}

echo -e "${GREEN}Building Docker image: ${IMAGE_NAME}:${TAG}${NC}"

echo -e "${YELLOW}Building Docker image...${NC}"
docker build -t "${IMAGE_NAME}:${TAG}" .

echo -e "${YELLOW}Testing Docker image...${NC}"
docker run --rm -d --name evilginx2-test -p 8080:8080 "${IMAGE_NAME}:${TAG}"

sleep 10

if curl -f http://localhost:8080/api/health > /dev/null 2>&1; then
    echo -e "${GREEN}Health check passed${NC}"
else
    echo -e "${RED}Health check failed${NC}"
    docker logs evilginx2-test
    docker stop evilginx2-test
    exit 1
fi

docker stop evilginx2-test

echo -e "${GREEN}Docker image built and tested successfully!${NC}"
echo -e "${GREEN}Image: ${IMAGE_NAME}:${TAG}${NC}"
