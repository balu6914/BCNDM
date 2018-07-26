# hyperledger-explorer-docker-image
This is one docker image of hyperledger explorer.

## Building your image
Go to the directory that has your Dockerfile and run the following command to build the Docker image. The -t flag lets you tag your image so it's easier to find later using the docker images command:
```sh
docker build -t datapace/hyperledger-explorer .
```

## Run the image
Running your image with -d runs the container in detached mode, leaving the container running in the background. The -p flag redirects a public port to a private port inside the container.

Run the image you previously built:
```
docker run -v /Users/datapace/blockchain/blockchain-explorer/config.json:/blockchain-explorer/config.json -v /Users/datapace/Downloads/fabric-samples-1.1.0/first-network/crypto-config:/first-network/crypto-config -p 8080:8080 -d datapace/hyperledger-explorer
```

- Mount your config.json to to /blockchain-explorer/config.json (path in container)
    - Change "fabric-path" to your fabric network path, example: "/home/user1/workspace/fabric-samples" for the following keys: "tls_cacerts", "key", "cert".
    - Final path for key "tls_cacerts" will be: "/first-network/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt".
- Modify config.json to update one of the channel
    - pg host, username, password details.

    ```
    "channel": "mychannel",
    "pg": {
		"host": "127.0.0.1",
		"port": "5432",
		"database": "fabricexplorer",
		"username": "hppoc",
		"passwd": "password"
	}
    ```
- Mount your crypto-config to {crypto-config} path in the container
    -  {crypto-config} should be matched to the value specified in /blockchain-explorer/config.json

**Tips**: you can use the following command to enter into the container to verify the configuration above.

> docker run -v /Users/datapace/blockchain/blockchain-explorer/config.json:/blockchain-explorer/config.json -v /Users/datapace/Downloads/fabric-samples-1.1.0/first-network/crypto-config:/first-network/crypto-config  -p 8080:8080   -it datapace/hyperledger-explorer bash

## Push Docker images to your namespace
Upload (push) the image to your namespace. Replace <my_namespace> with your namespace, and <image_repo> and <tag> with the repository and the tag of the image that you chose when you tagged the image.

```docker
docker push <my_namespace>/<image_repo>:<tag>
```

For example:
```sh
docker push datapace/hyperledger-explorer:latest
```
