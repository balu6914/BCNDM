#!/bin/bash

# This script will run docker composition and clean old ghost containers

# Stop all containers (if running)
docker-compose -f docker/fabric/docker-compose.yaml stop

# Remove old keystore (bug with fabric, container must be removed)
docker rm $(docker ps -a | grep -v "fabric-ca-pgsql" | awk '{print $1}')

# Run the network again
docker-compose -f docker/fabric/docker-compose.yml up
