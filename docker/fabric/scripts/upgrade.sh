#!/bin/bash

CHANNEL_PATH=./artifacts/datapacechannel.tx
CHANNEL_BLOCK=datapacechannel.block

ORDERER_URL=orderer.datapace.com:7050

CHANNEL_ID=datapacechannel

TOKEN_CHAIN_ID=token
TOKEN_CHAIN_PATH=github.com/chaincode/token
TOKEN_CHAIN_VER=$1
TOKEN_CHAIN_INIT_FN='{"Args":["init","{\"name\": \"Datapace Token\", \"symbol\": \"DPC\", \"decimals\": 8, \"totalSupply\": 100000000000000}"]}'

FEE_CHAIN_ID=fee
FEE_CHAIN_PATH=github.com/chaincode/system-fee
FEE_CHAIN_VER=$1
FEE_CHAIN_INIT_FN='{"Args":["init","{\"owner\": \"Admin@org1.datapace.com\", \"value\": 10000}"]}'

CONTRACTS_CHAIN_ID=contracts
CONTRACTS_CHAIN_PATH=github.com/chaincode/contracts
CONTRACTS_CHAIN_VER=$1
CONTRACTS_CHAIN_INIT_FN='{"Args":["init"]}'

ACCESS_CHAIN_ID=access
ACCESS_CHAIN_PATH=github.com/chaincode/access-requests
ACCESS_CHAIN_VER=$1
ACCESS_CHAIN_INIT_FN='{"Args":["init"]}'

TERMS_CHAIN_ID=terms
TERMS_CHAIN_PATH=github.com/chaincode/terms
TERMS_CHAIN_VER=$1
TERMS_CHAIN_INIT_FN='{"Args":["init"]}'

CERT_PATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/datapace.com/orderers/orderer.datapace.com/msp/tlscacerts/tlsca.datapace.com-cert.pem

if [[ "$2" == "token" ]]; then
  echo "Upgrading $2 to version $1"
  LOCATION=$PWD
  cd $GOPATH/src/github.com/chaincode/token
  # Install govendor tool
  go get -u github.com/kardianos/govendor
  # Fetch deps
  govendor sync
  cd $LOCATION
  peer chaincode install -n $TOKEN_CHAIN_ID -v $TOKEN_CHAIN_VER -p $TOKEN_CHAIN_PATH
  peer chaincode upgrade -o $ORDERER_URL -n $TOKEN_CHAIN_ID -v $TOKEN_CHAIN_VER -c "$TOKEN_CHAIN_INIT_FN" -C $CHANNEL_ID --tls --cafile $CERT_PATH
fi

if [[ "$2" == "fee" ]]; then
  echo "Upgrading $2 to version $1"
  LOCATION=$PWD
  cd $GOPATH/src/github.com/chaincode/system-fee
  # Install govendor tool
  go get -u github.com/kardianos/govendor
  # Fetch deps
  govendor sync
  cd $LOCATION
  peer chaincode install -n $FEE_CHAIN_ID -v $FEE_CHAIN_VER -p $FEE_CHAIN_PATH
  peer chaincode upgrade -o $ORDERER_URL -n $FEE_CHAIN_ID -v $FEE_CHAIN_VER -c "$FEE_CHAIN_INIT_FN" -C $CHANNEL_ID --tls --cafile $CERT_PATH
fi

if [[ "$2" == "contracts" ]]; then
  echo "Upgrading $2 to version $1"
  LOCATION=$PWD
  cd $GOPATH/src/github.com/chaincode/contracts
  # Install govendor tool
  go get -u github.com/kardianos/govendor
  # Fetch deps
  govendor sync
  cd $LOCATION
  # Install chaincode
  peer chaincode install -n $CONTRACTS_CHAIN_ID -v $CONTRACTS_CHAIN_VER -p $CONTRACTS_CHAIN_PATH
  peer chaincode upgrade -o $ORDERER_URL -n $CONTRACTS_CHAIN_ID -v $CONTRACTS_CHAIN_VER -c "$CONTRACTS_CHAIN_INIT_FN" -C $CHANNEL_ID --tls --cafile $CERT_PATH
fi

if [[ "$2" == "access-requests" ]]; then
  echo "Upgrading $2 to version $1"
  LOCATION=$PWD
  cd $GOPATH/src/github.com/chaincode/access-requests
  # Install govendor tool
  go get -u github.com/kardianos/govendor
  # Fetch deps
  govendor sync
  cd $LOCATION
  # Install chaincode
  peer chaincode install -n $ACCESS_CHAIN_ID -v $ACCESS_CHAIN_VER -p $ACCESS_CHAIN_PATH
  peer chaincode upgrade -o $ORDERER_URL -n $ACCESS_CHAIN_ID -v $ACCESS_CHAIN_VER -c "$ACCESS_CHAIN_INIT_FN" -C $CHANNEL_ID --tls --cafile $CERT_PATH
fi

if [[ "$2" == "terms" ]]; then
  echo "Upgrading $2 to version $1"
  LOCATION=$PWD
  cd $GOPATH/src/github.com/chaincode/terms
  # Install govendor tool
  go get -u github.com/kardianos/govendor
  # Fetch deps
  govendor sync
  cd $LOCATION
  peer chaincode install -n $TERMS_CHAIN_ID -v $TERMS_CHAIN_VER -p $TERMS_CHAIN_PATH
  peer chaincode upgrade -o $ORDERER_URL -n $TERMS_CHAIN_ID -v $TERMS_CHAIN_VER -c "$TERMS_CHAIN_INIT_FN" -C $CHANNEL_ID --tls --cafile $CERT_PATH
fi

echo "Done"

