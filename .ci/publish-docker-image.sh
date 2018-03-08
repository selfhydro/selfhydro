#!/bin/sh

set -e -x

#apt-get update
#
#apt-get install -y curl
#
#apt-get install \
#     apt-transport-https \
#     ca-certificates \
#     curl \
#     gnupg2 \
#     software-properties-common

#curl -fsSL https://download.docker.com/linux/$(. /etc/os-release; echo "$ID")/gpg | sudo apt-key add -
#
#sudo apt-key fingerprint 0EBFCD88
#
#add-apt-repository \
#   "deb [arch=amd64] https://download.docker.com/linux/$(. /etc/os-release; echo "$ID") \
#   $(lsb_release -cs) \
#   stable"
#
#apt-get install -y docker-ce
#
#curl -fsSL get.docker.com -o get-docker.sh
#sh get-docker.sh

docker run hello-world

TAG=$(cat selfhydro-release/tag)

docker build --tag selfhydro-release:${TAG} /selfhydro/.ci/docker/

docker login -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD}

docker push selfhydro-release:${TAG}