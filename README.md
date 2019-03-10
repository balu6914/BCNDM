# Datapace
[![Build Status](https://semaphoreci.com/api/v1/projects/f0520a21-dfcd-45ec-85bb-c14fda410c37/2174671/badge.svg)](https://semaphoreci.com/datapace/datapace)

IoT data marketplace based on blockchain.

## Install
```
cd $GOPATH/src
git@github.com:Datapace/datapace.git datapace
```

For read-only access (non-developer) you can use `https` and avoid setting necessary GitHub SSH keys:
```
cd $GOPATH/src
git clone https://github.com/Datapace/datapace.git datapace
```


## Set Dev Env

### Add SSH Key
Add your `ssh` key to GitHub.

If you used `https` to clone the repo, you must change the git remote to SSH in order to be able to send PRs:
```
git remote set-url origin git@github.com:datapace/datapace.git
```

### Configure git To Use SSH
Configure `git` to use `ssh` for GitLab, in order to enable
`dep` functionality, as explained [here](https://gist.github.com/shurcooL/6927554) -
otherwise it will break because it can not handle username/password prompt.

```
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

### Protobuf
Datapace comes with already generated all necesary `*.pb.go` files,
but in case you need to (re)gernerate them, command is:

```
make proto
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

## Deploy
Hyperledger Fabric network and Datapace dApp are deployed in the form of docker composition.

More information about tools necessary for deployment and deployment itself can
be found in [README](docker/README.md) in `docker` dir.

In order to build backend microservices (i.e. Datapace dApp) you will need to make dockers:

```
make dockers
```
### Logging
For production centralized logging run the following command:

```
make logging
```

This will start EFK looging dockers
This is a shorthand for:

```
docker-compose -f docker/efk/docker-compose.yml up
```


### Blockchain
Be sure that you generated all necessary crypto-material (as described in aforementioned `docker` [README](docker/README.md)):
```
make crypto
```

> N.B. Do not regenerate crypto material during the product life-cycle, otherwise you will loose access to all users
> (because Fabric CA will have new crypto keys and will not be in sync with MongoDB which holds the user metadata).

If all crypto-material is generated, then blockchain network can be started with:

```
make runbc
```

This is a shorthand for:

```
./docker/fabric/generate.sh
./docker/fabric/run.sh
```

For production centralized logging run the following command:
```
make runlogbc
```
Important note make sure the file `docker/fabric/.env` is containg the correct IP of the logging server

### Blockchain Explorer
Once blockchain network is up and running, you can start Hyperledger Explorer with:

```
make runexplorer
```

> Hyperledger Explorer docker image will be built previously when `make` (or `make dockers`) is called.


### Datapace dApp

#### Docker
In another terminal run:
```
make rundapp
```

which is a shorthand for:
```
docker-compose -f docker/docker-compose.yml up
```

For production centralized logging run the following command:

```
make runlogdapp
```

Important note make sure the file `docker/fabric/.env` is containg the correct IP of the logging server


UI will be available at [http://localhost:4200](http://localhost:4200)

#### Natively
To run locally compiled binaries (during dev process) instead of Docker containers
first assure that all `localhost` mappings are present in `/etc/hosts`:
```
drasko@Marx:~$ cat /etc/hosts
127.0.0.1	localhost
127.0.1.1	Marx
127.0.1.1	orderer.datapace.com
127.0.1.1	peer0.org1.datapace.com
127.0.1.1	ca.org1.datapace.com
...
```
This is needed because FabricSDK reads information from `config.yaml` and to avoid changing it drastically
it is easier just to create these mappings in `/etc/hosts`

To run previously compiled binaries execute:
```
make rundev
```

Happy hackin'!
