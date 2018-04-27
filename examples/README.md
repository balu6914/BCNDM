### Collection of examples

----

### Requirments

- Running Fabric network
- Hyperledger Fabric [cryptogen](http://hyperledger-fabric.readthedocs.io/en/release-1.1/commands/cryptogen-commands.html) utility.
- Generated Fabric key material using cryptogen utility.
- For `fabric-ca` examples running [fabric-ca](https://github.com/hyperledger/fabric-ca) service.


 ## Cryptogen
 In order to generate Fabric  Fabric key material and certificates, we must run cryptogen tool and pass crypto-config.yaml

Clone fabric repository
```
git clone https://github.com/hyperledger/fabric
```
Build cryptogen
```
make cryptogen
```
Navigate to `monetasa/examples`
```
cd `path_to_monetasa_project/monetasa/examples`
```

Generate certificates

```
~/path_to_your_fabric_folder/fabric/build/bin/cryptogen generate --config=./crypto-config.yaml

```
You will see auto generated `crypto-config` folder with all requiured keys. 
