#!/bin/sh

set -e -x

export GOPATH=$PWD

mkdir -p src/github.com/bchalk101/
cp -R ./selfhydro src/github.com/bchalk101/.

cd src/github.com/bchalk101/selfhydro
go get
go test -cover ./... | tee test_coverage.txt

mv test_coverage.txt $GOPATH/coverage-results/.

