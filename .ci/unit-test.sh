#!/bin/bash

set -e

export GOPATH=$PWD/selfhydro

go test
