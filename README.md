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
otherwise it will break beacuse it can not handle username/password prompt.

```
git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"
```

## Deploy
Hyperledger Fabric network and Datapace dApp are deployed in the form of docker composition.

More information can be found in [README](docker/README.md) in `docker` dir.

### Blockchain
If all crypto-material is generated (as described in aforementioned `docker` README), then blockchain network can be started with:

```
./docker/fabric/run.sh
```

### Datapace dApp
In another terminal run:

```
docker-compose -f docker/docker-compose.yml up
```

Happy hackin'!
