#! /bin/bash

set -e -x

apt-get update
apt-get install zip -y

VERSION=$(cat version/version)

cd ./selfhydro/db/
GOOS=linux go build -o dynamoDBTableCreater


zip selfhydro-state-db-release.zip dynamoDBTableCreater
mv selfhydro-state-db-release.zip ../../release
