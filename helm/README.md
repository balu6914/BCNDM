# Datapace Helm charts

This is a development, single node install of datapace on kubernetes!

For production use, please create production ready deployment.

Prior to deploying charts, please go to /tls dir for instructions on how to create certificates for ingress

To deploy chart from this directory follow these steps:

### Prerequisites

The Datapace uses crypto material to communicate to the blockchain.
This crypto material needs to be provided on a NFS server during the initial installation.
The directory setup and contents are the same as for the docker installation.
The nfs shared dir should have these subdirs providing crypto material:
- config  
- datapace-service-kvs  
- datapace-service-msp

In order for Datapace to communicate with blockchain, the bc->ip needs to be set to the IP of the host where the blockchain is running (deployed by using `make runbcdev` from the root directory).

The kubernetes install is done using the following commands:

- install helm dependencies

`helm dependency build`

- install helm on kubernetes

`helm install dpd . --set imageCredentials.username="yourusername" --set imageCredentials.password="yourpassword" `

Datapace images are kept in a protected docker registry, so you need to supply valid credentials in order to let kubernetes be able to pull images.

Please note that kubernetes exposes `ui.datapace.local` and `dproxy.datapace.local` hostnames.
In order to make your browser work with it, you need to edit /etc/hosts file on your machine so those hostnames resolve to the kubernetes machine.