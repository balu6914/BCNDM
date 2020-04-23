#!/bin/bash
source env.sh
ssh-keyscan -H $BCHOSTIP >> ~/.ssh/known_hosts
ssh-keyscan -H $DATAPACEHOSTIP >> ~/.ssh/known_hosts
scp /tmp/cryptogen ubuntu\@$BCHOSTIP:/home/ubuntu/
if [ $? -ne 0 ]; then
    echo "failed to copy cryptogen, is it in your /tmp dir?"
fi
scp /tmp/configtxgen ubuntu\@$BCHOSTIP:/home/ubuntu/
if [ $? -ne 0 ]; then
    echo "failed to copy configtxgen, is it in your /tmp dir?"
fi
echo "doing bc machine"
ssh ubuntu@$BCHOSTIP << EOF
  sudo chmod +x /home/ubuntu/cryptogen
  sudo chmod +x /home/ubuntu/configtxgen
  sudo cp /home/ubuntu/cryptogen /usr/local/bin/
  sudo cp /home/ubuntu/configtxgen /usr/local/bin/
  sudo DEBIAN_FRONTEND=noninteractive apt -qq update > /dev/null 2>&1
  sudo DEBIAN_FRONTEND=noninteractive apt -qq install docker-compose make -y > /dev/null 2>&1
  sudo usermod -a -G docker ubuntu
  git clone https://$GITHUBTOKEN:x-oauth-basic@github.com/Datapace/datapace.git
  cd datapace
  ./docker/fabric/generate.sh
  tar -czf helm.tar.gz helm/
EOF
echo "doing runbcdev"
ssh ubuntu@$BCHOSTIP << EOF
  cd datapace
  make runbcdev
  sudo tar -czf config.tar.gz config/
EOF
scp ubuntu\@$BCHOSTIP:/home/ubuntu/datapace/config.tar.gz /tmp
scp ubuntu\@$BCHOSTIP:/home/ubuntu/datapace/helm.tar.gz /tmp
scp /tmp/config.tar.gz ubuntu\@$DATAPACEHOSTIP:/home/ubuntu/
scp /tmp/helm.tar.gz ubuntu\@$DATAPACEHOSTIP:/home/ubuntu/
echo "doing kubernetes machine"
ssh ubuntu@$DATAPACEHOSTIP << EOF
  sudo mkdir -p /nfsshare
  sudo cp /home/ubuntu/config.tar.gz /nfsshare
  cd /nfsshare
  sudo tar -xf config.tar.gz
  sudo rm config.tar.gz
  sudo DEBIAN_FRONTEND=noninteractive apt -qq update > /dev/null 2>&1
  sudo DEBIAN_FRONTEND=noninteractive apt -qq install nfs-server -y > /dev/null 2>&1
  echo "/nfsshare *(rw,no_root_squash,no_subtree_check)" | sudo tee -a /etc/exports
  sudo systemctl restart nfs-server
  sudo snap install microk8s --classic > /dev/null 2>&1
  sudo snap install helm --classic > /dev/null 2>&1
  sudo /snap/bin/microk8s.status --wait-ready > /dev/null 2>&1
  sudo /snap/bin/microk8s.enable storage dns ingress > /dev/null 2>&1
  sudo usermod -a -G microk8s ubuntu
EOF
echo "doing helm"
ssh ubuntu@$DATAPACEHOSTIP << EOF
  mkdir -p /home/ubuntu/.kube
  /snap/bin/microk8s.config > /home/ubuntu/.kube/config
  /snap/bin/helm repo add bitnami https://charts.bitnami.com/bitnami
  tar -xf helm.tar.gz
  cd helm/datapace
  sed -i 's/ui.datapace.local$/$UIHOSTNAME/g' values.yaml
  sed -i 's/dproxy.datapace.local$/$DPROXYHOSTNAME/g' values.yaml
  sed -i 's/ip: 10.0.0.1$/ip: $BCHOSTIP/g' values.yaml
  sed -i 's/ip: 10.0.0.2$/ip: $DATAPACEHOSTIP/g' values.yaml
  /snap/bin/helm dependency build
  /snap/bin/helm install dpd . --set imageCredentials.username=$GITLABUSER --set imageCredentials.password=$GITLABTOKEN --kubeconfig /home/ubuntu/.kube/config
EOF
echo "All done"
