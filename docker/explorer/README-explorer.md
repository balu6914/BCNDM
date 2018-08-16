# hyperledger-explorer-docker-image
This is one docker image of hyperledger explorer.

## Building your image
Go to the directory that has your Dockerfile and run the following command to build the Docker image. The -t flag lets you tag your image so it's easier to find later using the docker images command:
```sh
docker build -t datapace/hyperledger-explorer .
```

## Run the image
Running your image with -d runs the container in detached mode, leaving the container running in the background. The -p flag redirects a public port to a private port inside the container.   
The image will run automaticly when starting fabric network.


