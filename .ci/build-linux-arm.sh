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
	CC=arm-linux-gnueabihf-gcc
)

apt-get update
apt-get install crossbuild-essential-armhf -y

export GOPATH=$PWD

mkdir -p src/github.com/bchalk101/

cp -R ./selfhydro src/github.com/bchalk101/.

cd src/github.com/bchalk101/selfhydro

go get
env ${GO_CROSS_ENV[@]} go build -o binary/selfhydro

ls -la

cd binary

#docker build .

scp selfhydro pi@10.1.1.6:/selfhydro/

#ssh pi@10.1.1.6 'docker kill selfhydro'
#ssh pi@10.1.1.6 'docker run -v /sys:/sys -v /selfhydro:/selfhydro selfhydro'

