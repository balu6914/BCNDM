#!/bin/bash

# exit when any command fails
set -e

CHAINCODE_FOLDER=helm/hlf-kube/chaincode
config_file=config/fabric/network.yaml

###
# Clean previous
###
if [ -d $CHAINCODE_FOLDER ]; then
  sudo rm -rf $CHAINCODE_FOLDER
fi

mkdir -p $CHAINCODE_FOLDER

chaincodes=$(yq ".network.chaincodes[].name" $config_file -c -r)

currentChaincodePtr=-1
for chaincode in $chaincodes; do
  currentChaincodePtr=$((currentChaincodePtr+1))
  dir=`yq .network.chaincodes["$currentChaincodePtr"].folder config/fabric/network.yaml | sed 's/"//g'`
  echo "creating ${CHAINCODE_FOLDER}/${chaincode}_1.0.tar.gz"
  tar --exclude='vendor' -C chaincode -cvzf ${CHAINCODE_FOLDER}/$chaincode.tar ${dir}/
done

