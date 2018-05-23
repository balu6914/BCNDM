#!/bin/bash
set -e
# This script generates crypto-config a fabric key material using fabric Cryptogen tool.
# Resource: http://hyperledger-fabric.readthedocs.io/en/release-1.1/commands/cryptogen-commands.html
# And generate network artifacts, channel config and genesis block using fabric configtxgen tool.
# Reference: http://hyperledger-fabric.readthedocs.io/en/release-1.1/commands/configtxgen.html

# NOTE: Required tools are cryptogen and configtxgen
# You can install it by running following command from fabric root directory

CRYPTO_CONF_PATH="examples/config/crypto-config.yaml"
CRYPTO_CONF_OUTPUT_PATH="examples/crypto-config"
GENESIS_BLOCK_PATH="../docker/artifacts/genesis.block"
GENESIS_BLOCK_PROFILE="MonetasaOrdererGenesis"
CH_OUTPUT_PATH="../docker/artifacts/myc.tx"
CH_PROFILE="MonetasaChannel"
CH_ID="myc"
MSG_DONE="Success! All done."
LOCATION=$PWD
BASE_COMPOSE_PATH='docker/base/docker-compose-base.yaml'


# Get hyperledger/fabric
FABRIC_PATH="$GOPATH/src/github.com/hyperledger/fabric"

if [ ! -d $FABRIC_PATH ]; then
  if [ ! -d "$GOPATH/src/github.com/hyperledger" ]; then
    mkdir $GOPATH/src/github.com/hyperledger
  fi
  cd $GOPATH/src/github.com/hyperledger
  git clone git@github.com:hyperledger/fabric.git && cd fabric
fi

cd $FABRIC_PATH
make cryptogen configtxgen

cd $LOCATION

$FABRIC_PATH/build/bin/cryptogen generate --config=$CRYPTO_CONF_PATH --output=$CRYPTO_CONF_OUTPUT_PATH

# configtxgen requires configtx.yaml on current location
cd examples/config

$FABRIC_PATH/build/bin/configtxgen  -outputBlock $GENESIS_BLOCK_PATH -profile  $GENESIS_BLOCK_PROFILE


$FABRIC_PATH/build/bin/configtxgen  -outputCreateChannelTx $CH_OUTPUT_PATH  -profile $CH_PROFILE -channelID $CH_ID


cd $LOCATION

# Change org ca pub key file name docker-compose env var
NEW_CERT_NAME=`find $LOCATION/examples/crypto-config/peerOrganizations/org1.monetasa.com/ca/ -type f -name "*_sk" -printf "%f\n"`
sed -i "s#FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/.*#FABRIC_CA_SERVER_CA_KEYFILE=/etc/hyperledger/fabric-ca-server-config/$NEW_CERT_NAME#" examples/docker/base/docker-compose-base.yaml
sed -i "s#FABRIC_CA_SERVER_TLS_KEYFILE=/etc/hyperledger/fabric-ca-server-config/.*#FABRIC_CA_SERVER_TLS_KEYFILE=/etc/hyperledger/fabric-ca-server-config/$NEW_CERT_NAME#" examples/docker/base/docker-compose-base.yaml

echo $MSG_DONE

exit 0
