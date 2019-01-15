#!/bin/bash
#
# Copyright (c) 2018 Datapace
#

###
# Launches all Datapace binaries (must be previously built).
#
# Expects that infrastructure (Fabric and Mongo DB) is already installed and running.
#
###

ROOT=$GOPATH/src/datapace
BUILD=$ROOT/build
UI=$ROOT/ui

# Kill all datapace-* stuff
function cleanup {
	pkill datapace
}

###
# Auth
###
DATAPACE_AUTH_HTTP_PORT=8080                \
DATAPACE_AUTH_GRPC_PORT=8081                \
DATAPACE_AUTH_DB_NAME=datapace-auth         \
DATAPACE_AUTH_DB_URL=localhost              \
DATAPACE_TRANSACTIONS_URL=localhost:8083    \
DATAPACE_AUTH_SECRET=EEE                    \
$BUILD/datapace-auth &

###
# Transactions
###
DATAPACE_TRANSACTIONS_HTTP_PORT=8082                		\
DATAPACE_TRANSACTIONS_GRPC_PORT=8083                		\
DATAPACE_TRANSACTIONS_DB_NAME=datapace-transactions 		\
DATAPACE_TRANSACTIONS_FABRIC_ADMIN=Admin@org1.datapace.com 	\
DATAPACE_TRANSACTIONS_DB_URL=localhost              		\
$BUILD/datapace-transactions &

###
# Streams
###
DATAPACE_STREAMS_HTTP_PORT=8084             \
DATAPACE_STREAMS_GRPC_PORT=8085             \
DATAPACE_STREAMS_DB_URL=localhost           \
DATAPACE_STREAMS_DB_NAME=datapace-streams   \
DATAPACE_AUTH_URL=localhost:8081            \
$BUILD/datapace-streams &

###
# Subscriptions
###
DATAPACE_SUBSCRIPTIONS_PORT=8086                        \
DATAPACE_SUBSCRIPTIONS_DB_URL=localhost                 \
DATAPACE_SUBSCRIPTIONS_DB_NAME=datapace-subscriptions   \
DATAPACE_AUTH_URL=localhost:8081                        \
$BUILD/datapace-subscriptions &

###
# UI
###
cd $UI
npm start &
cd -

trap cleanup EXIT

while : ; do sleep 1 ; done
