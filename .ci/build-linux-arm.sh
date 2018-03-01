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
	CC=arm-linux-gnueabihf-gcc
)

OUTPUT_DIR=$(pwd)/linux-arm-binary

export GOPATH=$(pwd)/gopath:$(pwd)/gopath/src/github.com/bchalk101/selfhydro/Godeps/_workspace
cd gopath/src/github.com/bchalk101/selfhydro

if which arm-linux-gnueabihf-gcc 1> /dev/null 2>&1; then
  env go get ./...
  env ${GO_CROSS_ENV[@]} go build -o "$OUTPUT_DIR/self-hydro"
  env ${GO_ENV[@]} go test ./... --cover -v
else
  echo "Cross-compile environment is not installed."
  exit 1
fi