#!/bin/bash

set -e

TAG=$(cat version/version)


ssh pi@water.local 'docker run --name selfhydro --restart=always bchalk/selfhydro:${TAG}'
