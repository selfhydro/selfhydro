#!/bin/bash

set -e -x

export GOPATH=$PWD

mkdir -p src/github.com/bchalk101/
cp -R ./selfhydro src/github.com/bchalk101/.

cd src/github.com/bchalk101/selfhydro
go get
go test -cover ./...

mv test_coverage.txt $GOPATH/coverage-results/.

