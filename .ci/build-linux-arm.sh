#!/bin/bash

set -e -x

GO_ENV=(
	CGO_ENABLED=1
)

GO_CROSS_ENV=(
	GOARCH=arm
	GOARM=7
	GOOS=linux
	CGO_ENABLED=1
)


export GOPATH=$PWD

mkdir -p src/github.com/bchalk101/
cp -R ./selfhydro src/github.com/bchalk101/.

OUTPUT_DIR=$(pwd)/linux-arm-binary

cd src/github.com/bchalk101/selfhydro

go get
env ${GO_CROSS_ENV[@]} go build -o "$OUTPUT_DIR/selfhydro"
