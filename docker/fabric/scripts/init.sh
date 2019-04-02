#!/bin/bash
set -e
# This script expedites the chaincode development process by automating the
# requisite channel create/join commands and chaincode deployment

CHANNEL_PATH=./artifacts/datapacechannel.tx
CHANNEL_BLOCK=datapacechannel.block

ORDERER_URL=orderer.datapace.com:7050

CHANNEL_ID=datapacechannel

TOKEN_CHAIN_ID=token
TOKEN_CHAIN_PATH=github.com/chaincode/token
TOKEN_CHAIN_VER=1.0
TOKEN_CHAIN_INIT_FN='{"Args":["init","{\"name\": \"Datapace Token\", \"symbol\": \"DPC\", \"decimals\": 8, \"totalSupply\": 100000000000000}"]}'

FEE_CHAIN_ID=fee
FEE_CHAIN_PATH=github.com/chaincode/system-fee
FEE_CHAIN_VER=1.0
FEE_CHAIN_INIT_FN='{"Args":["init","{\"owner\": \"Admin@org1.datapace.com\", \"value\": 10000}"]}'

CONTRACTS_CHAIN_ID=contracts
CONTRACTS_CHAIN_PATH=github.com/chaincode/contracts
CONTRACTS_CHAIN_VER=1.0
CONTRACTS_CHAIN_INIT_FN='{"Args":["init"]}'

ACCESS_CHAIN_ID=access
ACCESS_CHAIN_PATH=github.com/chaincode/access-requests
ACCESS_CHAIN_VER=1.0
ACCESS_CHAIN_INIT_FN='{"Args":["init"]}'

CERT_PATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/datapace.com/orderers/orderer.datapace.com/msp/tlscacerts/tlsca.datapace.com-cert.pem

LOCATION=$PWD

MSG_DONE="
#################################################################
       ########## Success! Network is ready. #########
#################################################################
"

# We use a pre-generated orderer.block and channel transaction artifact (datapace.tx),
# both of which are created using the configtxgen tool

# first we create the channel against the specified configuration in datapace.tx
# this call returns a channel configuration block - datapace.block - to the CLI container

peer channel create -o $ORDERER_URL  -c $CHANNEL_ID -f $CHANNEL_PATH  --tls --cafile $CERT_PATH

# now we will join the channel and start the chain with datapace.block serving as the
# channel's first block (i.e. the genesis block)
peer channel join -b $CHANNEL_BLOCK -o $ORDERER_URL

sleep 5

cd $GOPATH/src/github.com/chaincode/token
# Install govendor tool
go get -u github.com/kardianos/govendor

# Fetch deps
govendor sync

cd $LOCATION

# Install chaincode
peer chaincode install -n $TOKEN_CHAIN_ID -v $TOKEN_CHAIN_VER -p $TOKEN_CHAIN_PATH

cd $GOPATH/src/github.com/chaincode/system-fee
# Install govendor tool
go get -u github.com/kardianos/govendor

# Fetch deps
govendor sync

cd $LOCATION

# Install chaincode
peer chaincode install -n $FEE_CHAIN_ID -v $FEE_CHAIN_VER -p $FEE_CHAIN_PATH

cd $GOPATH/src/github.com/chaincode/contracts
# Install govendor tool
go get -u github.com/kardianos/govendor

# Fetch deps
govendor sync

cd $LOCATION

# Install chaincode
peer chaincode install -n $CONTRACTS_CHAIN_ID -v $CONTRACTS_CHAIN_VER -p $CONTRACTS_CHAIN_PATH

cd $GOPATH/src/github.com/chaincode/access-requests
# Install govendor tool
go get -u github.com/kardianos/govendor

# Fetch deps
govendor sync

cd $LOCATION

# Install chaincode
peer chaincode install -n $ACCESS_CHAIN_ID -v $ACCESS_CHAIN_VER -p $ACCESS_CHAIN_PATH

sleep 5

# Init/provision system with DPC
peer chaincode instantiate -o $ORDERER_URL -n $TOKEN_CHAIN_ID -v $TOKEN_CHAIN_VER -c "$TOKEN_CHAIN_INIT_FN" -C $CHANNEL_ID --tls --cafile $CERT_PATH
sleep 30

# Init/provision system with system fee
peer chaincode instantiate -o $ORDERER_URL -n $FEE_CHAIN_ID -v $FEE_CHAIN_VER -c "$FEE_CHAIN_INIT_FN" -C $CHANNEL_ID --tls --cafile $CERT_PATH
sleep 30

# Init/provision system with contracts
peer chaincode instantiate -o $ORDERER_URL -n $CONTRACTS_CHAIN_ID -v $CONTRACTS_CHAIN_VER -c "$CONTRACTS_CHAIN_INIT_FN" -C $CHANNEL_ID --tls --cafile $CERT_PATH
sleep 30

# Init/provision system with access control
peer chaincode instantiate -o $ORDERER_URL -n $ACCESS_CHAIN_ID -v $ACCESS_CHAIN_VER -c "$ACCESS_CHAIN_INIT_FN" -C $CHANNEL_ID --tls --cafile $CERT_PATH
sleep 30

echo $MSG_DONE

sleep 60000
exit 0
