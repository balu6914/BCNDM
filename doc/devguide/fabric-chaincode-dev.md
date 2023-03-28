# Developer Guide for chaincode development

This document will help you to run [Hyperledger Fabric](https://github.com/hyperledger/fabric) network in development, also how to install and init [chaincode](https://hyperledger-fabric.readthedocs.io/en/release-1.1/chaincode4ade.html) in order to
become available for use on network.

## Requirements
* curl
* Docker
* docker-compose version 1.14.0 or greater.
* Go language 1.9.x
* Node.js Runtime and NPM version 8.9.x or greater

Please Follow the detailed instractions (different OS) on [offcial hyperledger documentation](http://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html)

## Getting started
Clone the fabric sames repository from github

```git clone git@github.com:hyperledger/fabric-samples.git```

To run a development network we need only `chaincode-docker-devmode` fodler.

```cd chaincode-docker-devmode```

You will see docker-compose-simple.yaml file which is actually minimal fabric service composition required for development network.

We need four required hyperledger images.

[fabric-orderer](https://hub.docker.com/r/hyperledger/fabric-orderer/tags/)

[fabric-peer](https://hub.docker.com/r/hyperledger/fabric-peer/tags/)

[fabric-tools](https://hub.docker.com/r/hyperledger/fabric-tools/)

[fabric-ccenv](https://hub.docker.com/r/hyperledger/fabric-ccenv/tags/)

***NOTE:*** On docker hub, images are not tagged as latest, so you need to go on [dockerhub hyperledger repository](https://hub.docker.com/u/hyperledger/ ) find those 4 images (they are linked above) and check the latest version tag.
***Currently its 1.4 and we will use this tag.***

Next, Change the image tag in docker-compose-simple.yaml

Open the file and  for each service edit image tag like this:

| Before        | After
| ------------- |:-------------:|
| image: hyperledger/fabric-orderer      | image: hyperledger/fabric-peer:1.4 |
| image:  hyperledger/fabric-peer     | image: hyperledger/fabric-peer:1.4      |
| image: hyperledger/fabric-tools  | image: hyperledger/fabric-tools:1.4      |
| image: hyperledger/fabric-ccenv | image: hyperledger/fabric-ccenv:1.4 |


Save the file and exit.

***NOTE:*** We must do this because images are not tagged as latest on dockerhub and tag is missing in compose file (you just add it).

Lets run this docker compose and spin up our development network

```docker-compose -f docker-compose-simple.yaml up```

It will take a while to download all images




After some time you will see logs from docker container this means our network I s up and running. Lets check it


```docker ps```

Will show you running hyperledger containers (those 4 containers are required for local development).


```
CONTAINER ID        IMAGE                        COMMAND                  CREATED             STATUS              PORTS                                            NAMES
07521659e141        hyperledger/fabric-tools     "/bin/bash -c ./scri…"   3 minutes ago       Up 3 minutes                                                         cli
3f9e9e912741        hyperledger/fabric-ccenv     "/bin/bash -c 'sleep…"   3 minutes ago       Up 3 minutes                                                         chaincode
1762497a7b7a        hyperledger/fabric-peer      "peer node start --p…"   3 minutes ago       Up 3 minutes        0.0.0.0:7051->7051/tcp, 0.0.0.0:7053->7053/tcp   peer
723097fb44cc        hyperledger/fabric-orderer   "orderer"                3 minutes ago Up 3 minutes        0.0.0.0:7050->7050/tcp                           orderer
```


Our chaincode is awailable in ***parent chaincode folder*** (its mounted to all docker containers using volume option, so acutally our chaincode is available in all containers). You can easily test different chaincodes by adding them to the  ***chaincode parent directory*** and relaunching your network (e.g restarting docker composition). By default some chaincode examples are accessible in your chaincode container, they are comming as simple examples with fabric-samples project.

We can check this out by first navigating to ***parent chaincode folder (e.g in fabric-samples root)***
```cd ../```
Then accesing our chaincode source folder
```cd chaincode```

***IMPORTANT:*** You can put your new chaincode here (it will be available in chaincode docker container), that's how you should deploy your custom chaincode. In this guide we will use example chaincode that comes with fabric-sample repository.

### Interact with a chaincode ###
Open a new terminal window (this terminal will stay opened) and exec to running ***chaincode container*** using following command:

```docker exec -it chaincode bash```

From now on, we are inside a chaincode container created by fabric-ccenv image. Lets Compile our example chaincode:

```cd chaincode_folder```

We the deafault example so our chaincode_folder is `sacc` (when you copy you custom chaincode you will use your chaincode name).

Compile the chaincode

```go build```

Now run the chaincode by executing the following command:

```CORE_PEER_ADDRESS=peer:7052 CORE_CHAINCODE_ID_NAME=datapacec:0 ./sacc```


Great! From now on, our chaincode is running on the peer but its not yet installed, lets try to install it and actually interact with it using fabric cli tools.

Open a new terminal window, Terminal 2 (its required to stay opened) and exec to running fabric cli container using following command:

```docker exec -it cli bash```

From now on, we are in ***cli*** container created by fabric-tools image.
Install our chaincode executing following command:

```peer chaincode install -p chaincodedev/chaincode/sacc -n datapacec -v 0```

Next, we must initialize our chaincode

```peer chaincode instantiate -n datapacec -v 0 -c '{"Args":["a","10"]}' -C datapacechannel```

***NOTE:***  This chaincode is awailable in cli container because of ***parent chaincode folder*** (its mounted to docker container using volume option so its available in all containers). Its important to remember this, because when you write your own custom chaincode you just need to copy it to fabric-samples/chaincode and it will be automatically available (don't forget to restart docker composition when you copy new chaincode) in all docker containers, so you can install/test you chaincode.

Now issue an invoke the chaincode to change the value of “a” to “20” executing following command:

```peer chaincode invoke -n datapacec -c '{"Args":["set", "a", "20"]}' -C datapacechannel```


Finally, quer your chaincode. We should get  a value of 20 as a result, executing following command:


```peer chaincode query -n datapacec -c '{"Args":["query","a"]}' -C datapacechannel```


You will see response “Query Result: 20 ”

Congradulation, you just have a hyperledger network running in development mode with deployed smart contracts and we actually interact with one just now.


For more info checkout [Reference documentation](https://hyperledger-fabric.readthedocs.io/en/release-1.1/chaincode4ade.html)
