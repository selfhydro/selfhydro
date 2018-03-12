#!/bin/bash

set -e -x

echo $DEPLOY_KEY \
  | sed -e 's/\(KEY-----\)\s/\1\n/g; s/\s\(-----END\)/\n\1/g' \
  | sed -e '2s/\s\+/\n/g' > deploy_key

unset DEPLOY_KEY

set -ex

chmod 600 deploy_key

TAG=$(cat selfhydro-release/tag)

ssh -o StrictHostKeyChecking=no -i deploy_key pi@10.1.1.6 'docker rm -f selfhydro || true'
ssh -o StrictHostKeyChecking=no -i deploy_key pi@10.1.1.6 'docker run --name selfhydro --restart=always -v /sys:/sys -v /selfhydro:/selfhydro bchalk/selfhydro:${TAG}'
