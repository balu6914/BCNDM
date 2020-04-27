#!/bin/bash
. env.sh
echo "uninstalling datapace helm chart"
ssh ubuntu@$DATAPACEHOSTIP << EOF
  /snap/bin/microk8s.config > /home/ubuntu/kubeconfig
  cd helm/datapace
  /snap/bin/helm uninstall dpd --kubeconfig /home/ubuntu/kubeconfig
EOF
echo "datapace helm chart uninstalled"
