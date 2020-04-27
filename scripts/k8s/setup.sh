#!/bin/bash
. env.sh
ssh-keyscan -H $BCHOSTIP >> ~/.ssh/known_hosts
ssh-keyscan -H $DATAPACEHOSTIP >> ~/.ssh/known_hosts
. install-bc.sh
. run-bc.sh
. install-datapace.sh
. install-helm.sh
echo "All done"
