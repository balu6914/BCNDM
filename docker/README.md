# Fabric Network
---
Following those instractions you will be able to run dockerized Hyperledger Network fully provisioned, with Datapace token chaincoide deployed.

**NOTE:**  All commands are executed from project root.

### Cleaning Old Docker Images
**NOTE:**
This is optional and you should skip this step unless explicitly you want to remove all Hyperledger docker images. If needed, old Hyperledger docker images can be cleaned with:
```
docker rmi -f `docker images | grep hyperledger`
```

### Fabric Tools Installation
Make sure that you have all Fabric prerequisites and development tools installed:
```
go get github.com/hyperledger/fabric
cd $GOPATH/src/github.com/hyperledger/fabric
git checkout v1.1.0
make -j 16 cryptogen
make -j 16 configtxgen
cp build/bin/* $GOBIN
```

Make sure that `$GOBIN` is in your `PATH`.

Now you should be able to use tools globally, for example:

```
$ cryptogen version
cryptogen:
 Version: 1.1.0
 Go version: go1.10.1
 OS/Arch: linux/amd64
```

## Network and system provisioning

### Generate Crypto Material
From Datapace **project root** run generation script:
```
./docker/fabric/generate.sh
```

### Start/Restart the Network
If docker composition is running first stop it with `ctrl+c` then run:
```
./docker/fabric/run.sh
```
After few seconds you should see `Success Network` is ready message.

## Testing the Network
To confirm that network is working properly and that token chaincode is deployed and running on network
we will call chaincode from our docker `cli` container.

Enter docker `cli` container:

```
docker exec -it cli bash
```

Query chaincode:

```
peer chaincode query -C myc -n token -c '{"Args":["balance","{\"user\": \"testUser\"}"]}'
```

You should get response:
```
Query Result: {"user":"testUser","value":0}

```

This confirms that Datapace network is fully operational with running token chaincode.

---
# Datapace dApp
In order to run Datapace docker composition  you must first build docker images locally (they are not available online yet).

Execute following command from project root:
```
make dockers
```
Run docker composition:
```
docker-compose -f docker/docker-compose.yml up
```

Now go to `http://localhost` in your browser to see it in action.

**NOTE**: Datapace is available **ONLY** via HTTPS all HTTP requests will be redirected to HTTPS.
**IMPORTANT:** In development environment we are using selnf-signed SSL certificates so your browser will report you a SSL error fist time when you navigate to `http://localhost`, to solve this click on “ADVANCED” and then on “Proceed to <domain name> (unsafe)”.