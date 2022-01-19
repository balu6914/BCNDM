#!/bin/bash

# creates genesis block and certificates

# exit when any command fails
set -e

CRYPTO_CONF_PATH=config/fabric/crypto-config.yaml
CONFIG_TX_PATH=config/fabric/configtx.yaml
CRYPTO_CONF_DIR=config/crypto-config
config_file=config/fabric/network.yaml

FABRIC_CFG_PATH=config/fabric
EXPLORER_CFG_PATH=config/explorer
GENESIS_BLOCK_PATH=helm/hlf-kube/channel-artifacts/genesis.block
GENESIS_BLOCK_PROFILE=DatapaceOrdererGenesis
CH_OUTPUT_PATH=helm/hlf-kube/channel-artifacts/datapacechannel.tx
CH_PROFILE=DatapaceChannel
CHSYS_ID=datapacesyschannel
CH_ID=datapacechannel

###
# Clean previous
###
if [ -d $CRYPTO_CONF_DIR ]; then
  sudo rm -rf $CRYPTO_CONF_DIR
fi

if [ -d helm/hlf-kube/crypto-config ]; then
  sudo rm -rf helm/hlf-kube/crypto-config
fi

# generate certs
echo "-- creating certificates --"
cryptogen generate --config=$CRYPTO_CONF_PATH --output=$CRYPTO_CONF_DIR

# place holder empty folders for external peer orgs
externalPeerOrgs=$(yq '.ExternalPeerOrgs // empty' $CRYPTO_CONF_PATH -r -c)
if [ "$externalPeerOrgs" ]; then
    echo "-- creating empty folders for external peer orgs --"
    for peerOrgDomain in $(echo "$externalPeerOrgs" | jq -r '.[].Domain'); do
        echo "$peerOrgDomain"
        mkdir -p "$CRYPTO_CONF_DIR/peerOrganizations/$peerOrgDomain/msp"
    done
fi

# place holder empty folders for external orderer orgs
externalOrdererOrgs=$(yq '.ExternalOrdererOrgs // empty' $CRYPTO_CONF_PATH -r -c)
if [ "$externalOrdererOrgs" ]; then
    echo "-- creating empty folders for external orderer orgs --"
    for ordererOrg in $(echo "$externalOrdererOrgs" | jq -rc '.[]'); do
        # echo "$ordererOrg"
        ordererOrgDomain=$(echo "$ordererOrg" | jq -r '.Domain')
        echo "$ordererOrgDomain"
        mkdir -p "$CRYPTO_CONF_DIR/ordererOrganizations/$ordererOrgDomain/msp/tlscacerts"
        for ordererHostname in $(echo "$ordererOrg" | jq -r '.Specs[].Hostname'); do
            # echo "ordererHostname: $ordererHostname"
            mkdir -p "$CRYPTO_CONF_DIR/ordererOrganizations/$ordererOrgDomain/orderers/$ordererHostname.$ordererOrgDomain/tls"
        done
    done
fi

# generate genesis block
echo "-- creating genesis block --"
genesisProfile=$(yq '.network.genesisProfile' $config_file -r)
systemChannelID=$(yq '.network.systemChannelID' $config_file -r)
FABRIC_CFG_PATH=$FABRIC_CFG_PATH configtxgen -outputBlock $GENESIS_BLOCK_PATH -profile $GENESIS_BLOCK_PROFILE -channelID $CHSYS_ID
FABRIC_CFG_PATH=$FABRIC_CFG_PATH configtxgen -outputCreateChannelTx $CH_OUTPUT_PATH -profile $CH_PROFILE -channelID $CH_ID

cp -R config/crypto-config helm/hlf-kube
cp -R config/crypto-config helm/datapace
cp config/fabric/config.yaml helm/datapace

# prepare chaincodes
helm/scripts/prepare_chaincodes.sh

#copy configtx.yaml
cp $CONFIG_TX_PATH helm/hlf-kube

#copy explorer config files
cp -R $EXPLORER_CFG_PATH helm/hlf-kube

#copy channel artifacts
#mkdir -p helm/hlf-kube/channel-artifacts
#cp $GENESIS_BLOCK_PATH helm/hlf-kube/channel-artifacts
#cp $CH_OUTPUT_PATH helm/hlf-kube/channel-artifacts
