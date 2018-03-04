#!/bin/bash

set -e

export GOPATH=$PWD/selfhydro
cd ./selfhydro
go get
go test
