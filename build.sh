#!/bin/bash

echo -e "\n\n >> go mod download"
go mod download

docker compose down || echo -e "\n\n >> No service running"
docker rmi tuanloc1105/room-mate-finance:latest || echo -e "\n\n >> No images"

echo -e "\n\n >> go build -o ./go-app"
CGO_ENABLED=0 GOOS=linux go build -o ./go-app

echo -e "\n\n >> build images"
docker compose -f docker-compose.yml up -d room-mate-finance-go-service-2

echo -e "\n\n >> remove built go app"
rm -f ./go-app

echo -e "\n\n >> Done"
