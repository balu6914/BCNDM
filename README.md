# Monetasa
IoT data marketplace based on blockchain.

## Install
```
cd $GOPATH/src
git clone https://gitlab.com/drasko/monetasa
```
## Set Dev Env

### Add SSH Key
Add your `ssh` key to GitLab, then use ssh git remote:
```
git remote set-url origin git@gitlab.com:drasko/monetasa.git
```

### Configure git To Use SSH
Configure `git` to use `ssh` for GitLab, in order to enable
`dep` functionality, as explained [here](https://gist.github.com/shurcooL/6927554) -
otherwise it will break because it can not handle username/password prompt.

```
git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"
```

## Deploy
Hyperledger Fabric network and Datapace dApp are deployed in the form of docker composition.

More information can be found in [README](docker/README.md) in `docker` dir.

In order to build backend microservices (i.e. Datapace dApp) you will need to make dockers:

```
make proto
make dockers
```

Note that suggested version of `protobuf` is `3.6`. Outdated `protobuf` versions will
not work.

If recent `protoc` version is not provided by your system package manager, it can be installed manually from tarball:

```
wget https://github.com/google/protobuf/releases/download/v3.6.0/protoc-3.6.0-linux-x86_64.zip
unzip protoc-3.6.0-linux-x86_64.zip -d protoc3
sudo mv protoc3/bin/* /usr/local/bin/
sudo mv protoc3/include/* /usr/local/include/
export PATH=$PATH:/usr/local/bin/protoc
```

For the installation of Go `protobuf` tools and `gRPC` take a look [here](https://github.com/grpc/grpc-go#faq)

### Blockchain
If all crypto-material is generated (as described in aforementioned `docker` README), then blockchain network can be started with:

```
make runbc
```

This is a shorthand for:

```
./docker/fabric/generate.sh
./docker/fabric/run.sh
```

### Datapace dApp
In another terminal run:

```
docker-compose -f docker/docker-compose.yml up
```

To run locally compiled binaries (during dev process) instead of Docker containers execute:
```
make rundapp
```

UI will be available at [http://localhost:4200](http://localhost:4200)

Happy hackin'!
