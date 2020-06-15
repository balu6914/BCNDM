#!/bin/bash
. env.sh
echo "preparing kubernetes machine"
scp /tmp/config.tar.gz $DATAPACEHOSTUSER\@$DATAPACEHOSTIP:$DATAPACEHOSTHOME/
if [ $? -ne 0 ]; then
    echo "failed to copy config.tar.gz from /tmp, aborting"
    exit 1
fi
ssh $DATAPACEHOSTUSER@$DATAPACEHOSTIP << EOF
  sudo apt update
  sudo apt install -y snapd
  sudo mkdir -p /tmp/bc
  sudo cp $DATAPACEHOSTHOME/config.tar.gz /tmp/bc
  cd /tmp/bc
  sudo tar -xf config.tar.gz
  sudo rm config.tar.gz
  sudo snap install microk8s --classic
  sudo snap install helm --classic
  sudo /snap/bin/microk8s.status --wait-ready
  sudo /snap/bin/microk8s.enable storage dns ingress
  sudo usermod -a -G microk8s $DATAPACEHOSTUSER
EOF
echo "kubernetes machine done"
