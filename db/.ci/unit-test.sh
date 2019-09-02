#! /bin/bash

set -ex

export GOPATH=$PWD
cd ./selfhydro/db
go test -cover | tee test_coverage.txt

mv test_coverage.txt ../../coverage-results/.
