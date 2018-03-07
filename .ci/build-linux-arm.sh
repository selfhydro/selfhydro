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
env ${GO_CROSS_ENV[@]} go build -o release/selfhydro

ls -la

echo "v$(cat version/version)" > release/name
echo "v$(cat version/version)" > release/tag

cat > release/body <<EOF
Selfhydro compiled to be placed in docker image
EOF

