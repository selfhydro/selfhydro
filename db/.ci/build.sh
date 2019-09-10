#! /bin/bash

set -e -x

apt-get update
apt-get install zip -y

VERSION=$(cat version/version)

cd ./selfhydro/db/
GOOS=linux go build -v -ldflags '-d -s -w' -o dynamoDBTableCreater dynamoDBTableCreater.go


zip selfhydro-state-db-release.zip dynamoDBTableCreater
mv selfhydro-state-db-release.zip ../../release
