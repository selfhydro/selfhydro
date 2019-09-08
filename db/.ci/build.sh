#! /bin/bash

set -e -x

apt-get update
apt-get install zip -y

VERSION=$(cat version/version)

cd ./selfhydro/db/
GOOS=linux go build -o stateTableCreater


zip selfhydro-state-db-release.zip stateTableCreater
mv selfhydro-state-db-release.zip ../../release
