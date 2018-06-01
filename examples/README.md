# Collection of examples

## Prerequisits

### Cleaning Old Docker Images
If needed, old Hyperledger docker images can be cleaned with:
```
docker rmi -f `docker images | grep hyperledger`
```

### Fabric Tools Installation
Make sure that you have all Fabric prerequisites (docker images and dev tools) installed:
```
go get -u github.com/hyperledger/fabric
cd $GOPATH/github.com/hyperledger/fabric
git checkout v1.1.0
make -j 16 cryptogen
make -j 16 configtxgen
cp build/bin/* $GOBIN
```

Make sure that `$GOBIN` is in your `PATH`.

Now you should be able to use tools globally, for example:
```
drasko@Marx:~/go/src/github.com/hyperledger/fabric$ cryptogen version
cryptogen:
 Version: 1.1.0
 Go version: go1.10.1
 OS/Arch: linux/amd64
```

## Network and system provisioning

### Generate Crypto Material
From Datapace **project root** run generation script:
```
./examples/generate.sh
```

### Start/Restart the Network
If docker composition is running first stop it with `ctrl+c` then run:
```
./examples/run.sh
```
After few seconds you should see `Success Network` is ready message.

## Testing the Network
To confirm that network is working properly and that token chaincode is deployed and running on network
we will call chaincode from our docker `cli` container.

- Enter docker `cli` container:

```
docker exec -it cli bash
```

- Query chaincode:

```
peer chaincode query -C myc -n token -c '{"Args":["balance","{\"user\": \"testUser\"}"]}'
```

You should get response:
```
Query Result: {"user":"testUser","value":0}

```

This confirms that Datapace network is fully operational with running token chaincode.


## Testing Examples

You can test our examples
* fabric-ca - Creates new user in fabric-ca using fabric-sdk-go
* token-balance - Query our chaincode for user balance using  fabric-sdk-go

Navigate to example folder `cd EXAMPLE_FOLDER_NAME` and execute
```
go run main.go
```
