#!/bin/bash

set -e

TAG=$(cat version/version)


ssh pi@water.local 'docker run --name selfhydro --restart=always -v /sys:/sys -v /selfhydro:/selfhydro bchalk/selfhydro:${TAG}'
