#!/bin/bash
echo "preparing bc machine"
. env.sh
scp ./cryptogen $BCHOSTUSER\@$BCHOSTIP:$BCHOSTHOME/
if [ $? -ne 0 ]; then
    echo "failed to copy cryptogen, aborting, is it in this dir?"
    exit 1
fi
scp ./configtxgen $BCHOSTUSER\@$BCHOSTIP:$BCHOSTHOME/
if [ $? -ne 0 ]; then
    echo "failed to copy configtxgen, aborting, is it in this dir?"
    exit 1
fi
ssh $BCHOSTUSER@$BCHOSTIP << EOF
  sudo chmod +x $BCHOSTHOME/cryptogen
  sudo chmod +x $BCHOSTHOME/configtxgen
  sudo cp $BCHOSTHOME/cryptogen /usr/local/bin/
  sudo cp $BCHOSTHOME/configtxgen /usr/local/bin/
  sudo DEBIAN_FRONTEND=noninteractive apt -qq update
  sudo DEBIAN_FRONTEND=noninteractive apt -qq install docker-compose make -y
  sudo usermod -a -G docker $BCHOSTUSER
  git clone https://$GITHUBTOKEN:x-oauth-basic@github.com/Datapace/datapace.git
  cd datapace
  ./docker/fabric/generate.sh
  tar -czf helm.tar.gz helm/
  sudo tar -czf config.tar.gz config/
EOF
scp $BCHOSTUSER\@$BCHOSTIP:$BCHOSTHOME/datapace/config.tar.gz /tmp
scp $BCHOSTUSER\@$BCHOSTIP:$BCHOSTHOME/datapace/helm.tar.gz /tmp
echo "bc machine done"
