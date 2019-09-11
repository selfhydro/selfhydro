#! /bin/bash

set -e -x

apt-get update
apt-get install zip -y

VERSION=$(cat version/version)

cd ./selfhydro/db/
GOOS=linux GOARCH=amd64 go build -o dynamoDBTableCreater dynamoDBTableCreater.go
chmod +x dynamoDBTableCreater

zip selfhydro-state-db-release.zip dynamoDBTableCreater
mv selfhydro-state-db-release.zip ../../release
