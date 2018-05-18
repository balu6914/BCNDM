### Collection of examples

----
# Network and system provisioning
Monetasa network runs in docker

### Getting Started
`NOTE: Follow this step ONLY if you are running the network for first time (it means you don't have crypto material, genesis block, channel config etc...), otherwise jump to step 2. Running docker-compose`

1.From monetasa **project root** run generation script

```
./examples/generate.sh
```

After few seconds you should see Success message.


2.Run docker-compose from **project root**

```
docker-compose -f examples/docker/docker-compose.yaml up
```

After few seconds you should see Success Network is ready message.

### Testing network
To confirm that network is working properly and that token chaincode is deployed and running on network we will call chaincode from our docker `cli` container.

1.Enter docker `cli` container

```
docker exec -it cli bash
```

2.Query chaincode

```
peer chaincode query -C myc -n token -c '{"Args":["balance","{\"user\": \"testUser\"}"]}'

```

You should get response
```
Query Result: {"user":"testUser","value":0}

```

This confirm that Monetasa network is fully operational with running token chaincode.
