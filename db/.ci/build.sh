#! /bin/bash

set -e -x

VERSION=$(cat version/version)

cd ./selfhydro/db/
GOOS=linux go build -o stateTableCreater


zip selfhydro-state-db-release-$VERSION.zip stateTableCreater
mv selfhydro-state-db-release-$VERSION.zip ../../release
