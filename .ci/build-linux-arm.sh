#!/bin/bash

set -ex

GO_ENV=(
	CGO_ENABLED=1
)

GO_CROSS_ENV=(
	GOARCH=arm
	GOARM=7
	GOOS=linux
	CGO_ENABLED=1
	GOMAXPROCS=1
)

OUTPUT_DIR=$(pwd)/linux-arm-binary

export GOPATH=$(pwd)/gopath:$(pwd)/gopath/src/github.com/bchalk101/selfhydro/Godeps/_workspace
cd gopath/src/github.com/bchalk101/selfhydro

env  go get
env ${GO_CROSS_ENV[@]} go build -o "$OUTPUT_DIR/selfhydro"
env ${GO_ENV[@]} go test ./... --cover -v
