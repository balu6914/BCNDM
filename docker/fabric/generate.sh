#!/bin/bash

# This script generates crypto-config a fabric key material using fabric Cryptogen tool.
# (reference: http://hyperledger-fabric.readthedocs.io/en/release-1.1/commands/cryptogen-commands.html)
# and generate network artifacts, channel config and genesis block using fabric configtxgen tool.
# (reference: http://hyperledger-fabric.readthedocs.io/en/release-1.1/commands/configtxgen.html)

# NOTE: Required tools are cryptogen and configtxgen

CRYPTO_CONF_PATH=config/fabric/crypto-config.yaml
CRYPTO_CONF_DIR=config/crypto-config

FABRIC_CFG_PATH=config/fabric
GENESIS_BLOCK_PATH=docker/fabric/artifacts/genesis.block
GENESIS_BLOCK_PROFILE=DatapaceOrdererGenesis
CH_OUTPUT_PATH=docker/fabric/artifacts/datapacechannel.tx
CH_PROFILE=DatapaceChannel
CHSYS_ID=datapacesyschannel
CH_ID=datapacechannel

BASE_COMPOSE_FILE=docker/fabric/base/docker-compose-base.yaml

###
# Clean previous
###
if [ -d $CRYPTO_CONF_DIR ]; then
  sudo rm -rf $CRYPTO_CONF_DIR
fi

rm -rf /tmp/datapace-service-*

###
# Fabric keys
###
echo "### Generating Fabric keys"
cryptogen generate --config=$CRYPTO_CONF_PATH --output=$CRYPTO_CONF_DIR

###
# Network artifacts
###
echo "### Generating network artifacts"
FABRIC_CFG_PATH=$FABRIC_CFG_PATH configtxgen -outputBlock $GENESIS_BLOCK_PATH -profile $GENESIS_BLOCK_PROFILE -channelID $CHSYS_ID
FABRIC_CFG_PATH=$FABRIC_CFG_PATH configtxgen -outputCreateChannelTx $CH_OUTPUT_PATH -profile $CH_PROFILE -channelID $CH_ID

###
# Update docker-compose
###
NEW_CERT_NAME=`find $CRYPTO_CONF_DIR/peerOrganizations/org1.datapace.com/ca/ -type f -name "*_sk" | xargs basename`
sed -i -- "s#\(FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config\)/.*#\1/$NEW_CERT_NAME#" $BASE_COMPOSE_FILE
sed -i -- "s#\(FABRIC_CA_SERVER_TLS_KEYFILE=/etc/hyperledger/fabric-ca-server-config\)/.*#\1/$NEW_CERT_NAME#" $BASE_COMPOSE_FILE

echo "Success! All done."
