#!/bin/bash

set -e -x

GO_ENV=(
	CGO_ENABLED=1
)

GO_CROSS_ENV=(
    GOOS=linux
	GOARCH=arm
	GOARM=7
	CGO_ENABLED=1
	CC=arm-linux-gnueabi-gcc
)

sudo apt-get install binutils-arm-linux-gnueabi

export GOPATH=$PWD

mkdir -p src/github.com/bchalk101/

cp -R ./selfhydro src/github.com/bchalk101/.

cd src/github.com/bchalk101/selfhydro

go get
env ${GO_CROSS_ENV[@]} go build -o binary/selfhydro

ls -la

cd binary

scp selfhydro pi@water.local:/selfhydro/

ssh pi@water.local 'nohup sudo ./selfhydro/selfhydro &'

