#!/bin/bash
set -e
# This script expedites the chaincode development process by automating the
# requisite channel create/join commands and chaincode deployment

CHANNEL_PATH=./artifacts/myc.tx
CHANNEL_BLOCK=myc.block

ORDERER_URL=orderer.monetasa.com:7050

CHANNEL_ID=myc

TOKEN_CHAIN_ID=token
TOKEN_CHAIN_PATH=github.com/chaincode/token
TOKEN_CHAIN_VER=1.0
TOKEN_CHAIN_INIT_FN='{"Args":["init","{\"name\": \"Monetasa Token\", \"symbol\": \"TAS\", \"decimals\": 8, \"totalSupply\": 100000000000000}"]}'

FEE_CHAIN_ID=fee
FEE_CHAIN_PATH=github.com/chaincode/system-fee
FEE_CHAIN_VER=1.0
FEE_CHAIN_INIT_FN='{"Args":["init","{\"owner\": \"Admin@org1.monetasa.com\", \"value\": 10000}"]}'

CERT_PATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/monetasa.com/orderers/orderer.monetasa.com/msp/tlscacerts/tlsca.monetasa.com-cert.pem

LOCATION=$PWD

MSG_DONE="
#################################################################
       ########## Success! Network is ready. #########
#################################################################
"

# We use a pre-generated orderer.block and channel transaction artifact (myc.tx),
# both of which are created using the configtxgen tool

# first we create the channel against the specified configuration in myc.tx
# this call returns a channel configuration block - myc.block - to the CLI container

peer channel create -o $ORDERER_URL  -c $CHANNEL_ID -f $CHANNEL_PATH  --tls --cafile $CERT_PATH

# now we will join the channel and start the chain with myc.block serving as the
# channel's first block (i.e. the genesis block)
peer channel join -b $CHANNEL_BLOCK -o $ORDERER_URL

sleep 5

cd $GOPATH/src/github.com/chaincode/token
# Install govendor tool
go get -u github.com/kardianos/govendor

govendor init

# Fetch deps
govendor fetch github.com/hyperledger/fabric/protos/msp

cd $LOCATION

# Install chaincode
peer chaincode install -n $TOKEN_CHAIN_ID -v $TOKEN_CHAIN_VER -p $TOKEN_CHAIN_PATH

cd $GOPATH/src/github.com/chaincode/system-fee
# Install govendor tool
go get -u github.com/kardianos/govendor

govendor init

# Fetch deps
govendor fetch github.com/hyperledger/fabric/protos/msp

cd $LOCATION

# Install chaincode
peer chaincode install -n $FEE_CHAIN_ID -v $FEE_CHAIN_VER -p $FEE_CHAIN_PATH

sleep 5

# Init/provision system with TAS
peer chaincode instantiate -o $ORDERER_URL -n $TOKEN_CHAIN_ID -v $TOKEN_CHAIN_VER -c "$TOKEN_CHAIN_INIT_FN" -C $CHANNEL_ID --tls --cafile $CERT_PATH

# Init/provision system with system fee
peer chaincode instantiate -o $ORDERER_URL -n $FEE_CHAIN_ID -v $FEE_CHAIN_VER -c "$FEE_CHAIN_INIT_FN" -C $CHANNEL_ID --tls --cafile $CERT_PATH

sleep 5

echo $MSG_DONE

sleep 60000
exit 0
