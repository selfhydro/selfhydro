#!/bin/bash

set -e -x

TAG=$(cat selfhydro-release/tag)

ssh -o StrictHostKeyChecking=no pi@10.1.1.8 'docker rm -f selfhydro'
ssh -o StrictHostKeyChecking=no pi@10.1.1.8 'docker run --name selfhydro --restart=always -v /sys:/sys -v /selfhydro:/selfhydro bchalk/selfhydro:${TAG}'
