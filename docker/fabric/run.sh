#!/bin/bash

# This script will run docker composition and clean old ghost containers

docker-compose -f docker/fabric/docker-compose.yaml stop
# Remove CLI stoped container (bug with fabric tools image, container must be removed)
docker rm cli
# Run the network again
docker-compose -f docker/fabric/docker-compose.yaml up
