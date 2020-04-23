# K8S datapace installation script


The setup consists of three machines:
1. **LOCAL** machine where this file and the `setup.sh`, `env.sh` reside and from which the installation will be initiated.
2. **BC** machine where blockchain should be installed
3. **K8S** machine where datapace should be installed

## Prerequisites:

- on **LOCAL** machine **cryptogen** and **configtxgen** binaries should be put in the /tmp dir
- on **BC** machine ports 7050, 7051, 7053 and 7054 should be open  to the outside world
- on **K8S** machine ports 80 and 443 should be open to the outside world

This script should be run under assumption that:
- the servers are default Ubuntu installations
- you can ssh to the server from LOCAL machine by using your private key.
That is, when you execute ssh ubuntu@BC or ubuntu@K8S, you should be logged in without password prompt.

## Installation

Prior to running the script, change the environment variables in env.sh

- BCHOSTIP

this is the IP address of the BC machine

- DATAPACEHOSTIP

this is the IP address of the K8S machine

- GITHUBTOKEN

this is your **github** token you need to generate and give it sufficient privileges so git clone of datapace repo can be done

- GITLABUSER

this is your **gitlab** username of user which can access gitlab registry with datapace docker images

- GITLABTOKEN

this is your **gitlab** token you need to generate and give it sufficient privileges so docker pull can be executed with it

- UIHOSTNAME

this is the host which you will enter in your browser to access datapace ui

- DPROXYHOSTNAME

this is dproxy hostname which will deliver proxied files

Once you set up all the variables in env.sh, execute `sh setup.sh` from the same directory where the file resides and observe the progress.
When the script finishes it will output the 'All done' message.
You can then start your browser and go do the UIHOSTNAME that you previously specified, where you should be presented with the Datapace UI.
