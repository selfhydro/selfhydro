#!/bin/bash

set -e -x

echo $DEPLOY_KEY \
  | sed -e 's/\(KEY-----\)\s/\1\n/g; s/\s\(-----END\)/\n\1/g' \
  | sed -e '2s/\s\+/\n/g' > deploy_key

unset DEPLOY_KEY

set -ex

chmod 600 deploy_key

TAG=$(cat selfhydro-release/tag)

ssh -o StrictHostKeyChecking=no -i deploy_key pi@water.local << EOF
sudo systemctl stop selfhydro.service
EOF

scp -o StrictHostKeyChecking=no -i deploy_key selfhydro-release/selfhydro pi@10.1.1.3:/selfhydro/
ssh -o StrictHostKeyChecking=no -i deploy_key pi@water.local << EOF
chmod +x /selfhydro/selfhydro &\
sudo systemctl start selfhydro.service
EOF

#ssh -o StrictHostKeyChecking=no  pi@10.1.1.2 'docker rm -f selfhydro || true '
#ssh -o StrictHostKeyChecking=no  pi@10.1.1.2  "docker run --name selfhydro --restart=always --privileged=true -v /selfhydro:/selfhydro bchalk/selfhydro-release:$TAG"
