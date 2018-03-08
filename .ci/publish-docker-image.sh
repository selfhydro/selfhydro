#!/bin/bash

apt-get update
apt-get install -y curl

apt-get install docker-ce
docker run hello-world

TAG=$(cat selfhydro-release/tag)

docker build --tag selfhydro-release:${TAG} /selfhydro/.ci/docker/

docker login -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD}

docker push selfhydro-release:${TAG}