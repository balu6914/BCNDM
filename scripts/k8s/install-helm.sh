#!/bin/bash
. env.sh
scp /tmp/helm.tar.gz ubuntu\@$DATAPACEHOSTIP:/home/ubuntu/
if [ $? -ne 0 ]; then
    echo "failed to copy helm.tar.gz from /tmp, aborting"
    exit 1
fi
echo "installing datapace helm chart"
ssh ubuntu@$DATAPACEHOSTIP << EOF
  /snap/bin/microk8s.config > /home/ubuntu/kubeconfig
  /snap/bin/helm repo add bitnami https://charts.bitnami.com/bitnami
  tar -xf helm.tar.gz
  cd helm/datapace
  sed -i 's/ui.datapace.local$/$UIHOSTNAME/g' values.yaml
  sed -i 's/dproxy.datapace.local$/$DPROXYHOSTNAME/g' values.yaml
  sed -i 's/ip: 10.0.0.1$/ip: $BCHOSTIP/g' values.yaml
  /snap/bin/helm dependency build
  /snap/bin/helm install dpd . --set imageCredentials.username=$GITLABUSER --set imageCredentials.password=$GITLABTOKEN --kubeconfig /home/ubuntu/kubeconfig
EOF
echo "datapace helm chart done"
