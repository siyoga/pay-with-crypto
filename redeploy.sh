#!/bin/sh
docker-compose -f docker-compose.prod.yml --env-file prod.env down
git pull

go mod tidy
go build
docker-compose -f docker-compose.prod.yml --env-file prod.env build
docker-compose -f docker-compose.prod.yml --env-file prod.env up
