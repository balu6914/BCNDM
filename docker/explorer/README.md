# Hyperledger Explorer
Docker image of hyperledger explorer.

## Build
From the Datapace root dir execute:
```sh
docker build -f docker/explorer/Dockerfile -t datapace/hyperledger-explorer .
```

> `-t` flag lets you tag your image so it's easier to find later using the `docker images` command.

## Deploy  
Started as an addon to Hyperledger Fabric `docker-compose`:
```
docker-compose -f docker/fabric/docker-compose.yaml -f docker/explorer/docker-compose.yaml up
```

Explorer UI will be available at [http://localhost:9090](http://localhost:9090)
