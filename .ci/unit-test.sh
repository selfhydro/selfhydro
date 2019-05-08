#!/bin/bash

set -e

export GOPATH=$PWD

mkdir -p src/github.com/selfhydro/
cp -R ./selfhydro src/github.com/selfhydro/.

cd src/github.com/selfhydro/selfhydro
go get
go test -cover ./... | tee test_coverage.txt

mv test_coverage.txt $GOPATH/coverage-results/.

