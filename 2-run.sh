#!/bin/bash

clear

GO_VERSION=$(go version)
GO_PATH=$(where go || which go)
GO_ROOT=$(go env GOPATH)
GO_CACHE=$(go env GOCACHE)

echo -e " \n\n >> Go version: $GO_VERSION\n >> Go path: $GO_PATH\n >> Go root: $GO_ROOT\n >> Go cache: $GO_CACHE \n\n"

echo -e " \n\n >> Clean build cache \n\n"

go clean -cache

echo -e " \n\n >> Build service \n\n"

go get -u
go mod tidy
# go build

echo -e " \n\n >> Start running service \n\n"

go run main.go
