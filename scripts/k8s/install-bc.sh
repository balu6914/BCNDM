#!/bin/bash
echo "preparing bc machine"
. env.sh
scp ./cryptogen ubuntu\@$BCHOSTIP:/home/ubuntu/
if [ $? -ne 0 ]; then
    echo "failed to copy cryptogen, aborting, is it in this dir?"
    exit 1
fi
scp ./configtxgen ubuntu\@$BCHOSTIP:/home/ubuntu/
if [ $? -ne 0 ]; then
    echo "failed to copy configtxgen, aborting, is it in this dir?"
    exit 1
fi
ssh ubuntu@$BCHOSTIP << EOF
  sudo chmod +x /home/ubuntu/cryptogen
  sudo chmod +x /home/ubuntu/configtxgen
  sudo cp /home/ubuntu/cryptogen /usr/local/bin/
  sudo cp /home/ubuntu/configtxgen /usr/local/bin/
  sudo DEBIAN_FRONTEND=noninteractive apt -qq update
  sudo DEBIAN_FRONTEND=noninteractive apt -qq install docker-compose make -y
  sudo usermod -a -G docker ubuntu
  git clone https://$GITHUBTOKEN:x-oauth-basic@github.com/Datapace/datapace.git
  cd datapace
  ./docker/fabric/generate.sh
  tar -czf helm.tar.gz helm/
  sudo tar -czf config.tar.gz config/
EOF
scp ubuntu\@$BCHOSTIP:/home/ubuntu/datapace/config.tar.gz /tmp
scp ubuntu\@$BCHOSTIP:/home/ubuntu/datapace/helm.tar.gz /tmp
echo "bc machine done"
