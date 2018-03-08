#!/bin/bash

TAG=$(cat version/version)

docker build --tag selfhydro-release:${TAG} /selfhydro/.ci/docker/

docker login -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD}

docker push selfhydro-release:${TAG}