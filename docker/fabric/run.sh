#!/bin/bash

# This script will run docker composition and clean old ghost containers

docker-compose -f docker/fabric/docker-compose.yaml stop
# Remove CLI stoped container (bug with fabric tools image, container must be removed)
# Remove old keystore
#rm -rf /tmp/monetasa-service-*
#docker rm -f $(docker ps -a -q | grep -v "fabric-ca-pgsql")
docker rm $(docker ps -a | grep -v "fabric-ca-pgsql" | awk '{print $1}')
# Run the network again
docker-compose -f docker/fabric/docker-compose.yaml up
