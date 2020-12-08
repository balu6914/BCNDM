#!/bin/bash
set -e
# This script expedites the chaincode development process by automating the
# requisite channel create/join commands and chaincode deployment

CHANNEL_PATH=./artifacts/datapacechannel.tx
CHANNEL_BLOCK=datapacechannel.block

ORDERER_URL=orderer.datapace.com:7050

CHANNEL_ID=datapacechannel
ORG1MSP=Org1MSP

PEER_COUNT=(0 1)
CHAINCODES=("token" "contracts" "access" "terms")

TOKEN_CHAIN_ID=token
TOKEN_CHAIN_PATH=chaincode/token
TOKEN_CHAIN_VER=1.0
TOKEN_CHAIN_INIT_FN='{"Args":["init","{\"name\": \"Datapace Token\", \"symbol\": \"DPC\", \"decimals\": 8, \"totalSupply\": 100000000000000}"]}'

FEE_CHAIN_ID=fee
FEE_CHAIN_PATH=chaincode/system-fee
FEE_CHAIN_VER=1.0
FEE_CHAIN_INIT_FN='{"Args":["init","{\"owner\": \"Admin@org1.datapace.com\", \"value\": 10000}"]}'

CONTRACTS_CHAIN_ID=contracts
CONTRACTS_CHAIN_PATH=chaincode/contracts
CONTRACTS_CHAIN_VER=1.0
CONTRACTS_CHAIN_INIT_FN='{"Args":["init"]}'

ACCESS_CHAIN_ID=access
ACCESS_CHAIN_PATH=chaincode/access-requests
ACCESS_CHAIN_VER=1.0
ACCESS_CHAIN_INIT_FN='{"Args":["init"]}'

TERMS_CHAIN_ID=terms
TERMS_CHAIN_PATH=chaincode/terms
TERMS_CHAIN_VER=1.0
TERMS_CHAIN_INIT_FN='{"Args":["init"]}'
# Add chaincodes to the list, so we can itterate during installation process (every chaincode on every organization peer)
CHAINCODES_LIST=($TOKEN_CHAIN_ID $FEE_CHAIN_ID $CONTRACTS_CHAIN_ID $ACCESS_CHAIN_ID $TERMS_CHAIN_ID)

CERT_PATH=/opt/gopath/crypto/ordererOrganizations/datapace.com/orderers/orderer.datapace.com/msp/tlscacerts/tlsca.datapace.com-cert.pem
ORG_ROOT_CA_PATH=/opt/gopath/crypto/peerOrganizations/org1.datapace.com/peers/peer0.org1.datapace.com/tls/ca.crt
# Add peers from each organization, so we can chain peer connection params where multiple peers are required (approval, commit) 
ORG1_PEER0_ADDRESS=peer0.org1.datapace.com:7051
ORG1_PEER1_ADDRESS=peer1.org1.datapace.com:7051
PEER_CONN_PARMS="--peerAddresses $ORG1_PEER0_ADDRESS --peerAddresses $ORG1_PEER1_ADDRESS  --tlsRootCertFiles $ORG_ROOT_CA_PATH  --tlsRootCertFiles $ORG_ROOT_CA_PATH"

MSG_DONE="
#################################################################
       ########## Success! Network is ready. #########
#################################################################
"

# We use a pre-generated orderer.block and channel transaction artifact (datapacechannel.tx),
# both of which are created using the configtxgen tool

# First we create the channel against the specified configuration in datapacechannel.tx
# this call returns a channel configuration block - datapacechannel.block - to the CLI container
peer channel create -o $ORDERER_URL  -c $CHANNEL_ID -f $CHANNEL_PATH  --tls --cafile $CERT_PATH

# Join every peer to the channel
CORE_PEER_ADDRESS=$ORG1_PEER0_ADDRESS peer channel join -b $CHANNEL_BLOCK -o $ORDERER_URL
CORE_PEER_ADDRESS=$ORG1_PEER1_ADDRESS peer channel join -b $CHANNEL_BLOCK -o $ORDERER_URL

# Now we will loop over all peers from Datapce Org and install all chaincodes, approve them and commit approval. Then start the chain with datapacechannel.block serving as the
# channel's first block (i.e. the genesis block)
for peer in ${PEER_COUNT[@]}
do
       # Install, approve and commit to channel every chaincode in a list
       for chain in ${CHAINCODES_LIST[@]}
       do
              CURRENT_CHAIN=${chain^^}
              eval CH_ID='$'"$CURRENT_CHAIN"_CHAIN_ID
              eval CH_PATH='$'"$CURRENT_CHAIN"_CHAIN_PATH
              eval CH_VER='$'"$CURRENT_CHAIN"_CHAIN_VER
              eval CH_INIT_FN='$'"$CURRENT_CHAIN"_CHAIN_INIT_FN

              # Package chaincode
              peer lifecycle chaincode package $CH_ID'_'$CH_VER.tar.gz --path $CH_PATH --lang golang --label $CH_ID'_'$CH_VER

              # Install chaincode
              CORE_PEER_ADDRESS=peer$peer.org1.datapace.com:7051 peer lifecycle chaincode install $CH_ID'_'$CH_VER'.tar.gz'
              # Approve a chaincode definition as Datapace org and Commit the chaincode definition to the channel.
               if [ "$peer" -eq "0" ]; then
                     # Get package ID from package label. We need this for the next approval step
                     CORE_PEER_ADDRESS=peer$peer.org1.datapace.com:7051 peer lifecycle chaincode queryinstalled >&log.txt
                     PACKAGE_ID=$(sed -n "/${CH_ID}_${CH_VER}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)

                     # Approve a chaincode definition
                     CORE_PEER_ADDRESS=peer$peer.org1.datapace.com:7051 peer lifecycle chaincode approveformyorg -o $ORDERER_URL --ordererTLSHostnameOverride orderer.datapace.com --channelID $CHANNEL_ID --name $CH_ID --version $CH_VER --init-required --waitForEvent --package-id $PACKAGE_ID --sequence 1 --tls --cafile $CERT_PATH
                     # Wait for Org approvals
                     sleep 5
                     ORGAPPROVAL=$(peer lifecycle chaincode checkcommitreadiness --channelID ${CHANNEL_ID} --name ${CH_ID} --version ${CH_VER} --init-required --sequence 1 --tls --cafile ${ORDERER_CA} --output json | jq -r '.approvals') 
                     # We have only Datapace org for now, so majority policy is satisfied with one approval.
                     ORG1APPROVAL=$(echo $ORGAPPROVAL | jq -r --arg key "${ORG1MSP}" '.[$key]')

                     if [ "$ORG1APPROVAL" = true ] ; then
                            # Commit new chaincode definition on channel peers.
                            # NOTE: We need root CA files for each organization, multiple organization requires multiple --tlsRootCertFiles flags
                            CORE_PEER_ADDRESS=peer$peer.org1.datapace.com:7051 peer lifecycle chaincode commit -o $ORDERER_URL --ordererTLSHostnameOverride orderer.datapace.com --channelID $CHANNEL_ID --name $CH_ID --version $CH_VER --init-required --sequence 1 --tls --cafile $CERT_PATH $PEER_CONN_PARMS
                     fi
              fi
       done
done

sleep 5
# Invoke all chaincodes with init function - This process is done only once and init is required by organization policy (chaincode wont be accesible without init).
peer chaincode invoke -o $ORDERER_URL -n $TOKEN_CHAIN_ID --isInit -c "$TOKEN_CHAIN_INIT_FN" -C $CHANNEL_ID --tls --cafile $CERT_PATH $PEER_CONN_PARMS 
sleep 5

# Init/provision system with system fee
peer chaincode invoke -o $ORDERER_URL -n $FEE_CHAIN_ID --isInit -c "$FEE_CHAIN_INIT_FN" -C $CHANNEL_ID --tls --cafile $CERT_PATH $PEER_CONN_PARMS
sleep 5

# Init/provision system with contracts
peer chaincode invoke -o $ORDERER_URL -n $CONTRACTS_CHAIN_ID --isInit -c "$CONTRACTS_CHAIN_INIT_FN" -C $CHANNEL_ID --tls --cafile $CERT_PATH $PEER_CONN_PARMS
sleep 5

# Init/provision system with access control
peer chaincode invoke -o $ORDERER_URL -n $ACCESS_CHAIN_ID --isInit -c "$ACCESS_CHAIN_INIT_FN" -C $CHANNEL_ID --tls --cafile $CERT_PATH $PEER_CONN_PARMS
sleep 5

# Init/provision system with terms
peer chaincode invoke -o $ORDERER_URL -n $TERMS_CHAIN_ID --isInit -c "$TERMS_CHAIN_INIT_FN" -C $CHANNEL_ID --tls --cafile $CERT_PATH $PEER_CONN_PARMS
sleep 5

echo $MSG_DONE

exit 0
