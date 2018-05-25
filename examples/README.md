# Collection of examples

## Prerequisits
### Fabric Installation
Make sure that you have all Fabric prerequisits (docker images and dev tools) installed:
```
curl -sSL https://goo.gl/kFFqh5 | bash -s 1.1.0
```

### Cleaning Old Docker Images
If needed, old Hyperledger docker images can be cleaned with:
```
docker rmi -f `docker images | grep hyperledger`
```

## Network and system provisioning

### Generate Crypto Material
From Datapace **project root** run generation script:
```
./examples/generate.sh
```

### Start the Network
```
docker-compose -f examples/docker/docker-compose.yaml up
```

After few seconds you should see `Success Network` is ready message.


N.B. sometimes it is needed to clean old docker containers:
```
docker rm `docker ps -a | grep hyperledger | awk '{print $1}'`
```

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
