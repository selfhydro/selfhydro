#!/bin/bash

set -e

TAG=$(cat selfhydro-release/tag)


ssh pi@water.local 'docker run --name selfhydro --restart=always -v /sys:/sys -v /selfhydro:/selfhydro bchalk/selfhydro:${TAG}'
