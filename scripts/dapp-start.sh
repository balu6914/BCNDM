#!/bin/bash
#
# Copyright (c) 2018 Datapace
#

###
# Launches all Monetasa binaries (must be previously built).
#
# Expects that infrastructure (Fabric and Mongo DB) is already installed and running.
#
###

ROOT=$GOPATH/src/monetasa
BUILD=$ROOT/build
UI=$ROOT/ui

# Kill all monetasa-* stuff
function cleanup {
	pkill monetasa
}

###
# Auth
###
MONETASA_AUTH_HTTP_PORT=8080                \
MONETASA_AUTH_GRPC_PORT=8081                \
MONETASA_AUTH_DB_NAME=monetasa-auth         \
MONETASA_AUTH_DB_URL=localhost              \
MONETASA_TRANSACTIONS_URL=localhost:8083    \
MONETASA_AUTH_SECRET=EEE                    \
$BUILD/monetasa-auth &

###
# Transactions
###
MONETASA_TRANSACTIONS_HTTP_PORT=8082                \
MONETASA_TRANSACTIONS_GRPC_PORT=8083                \
MONETASA_TRANSACTIONS_DB_NAME=monetasa-transactions \
MONETASA_TRANSACTIONS_DB_URL=localhost              \
$BUILD/monetasa-transactions &

###
# Streams
###
MONETASA_STREAMS_HTTP_PORT=8084             \
MONETASA_STREAMS_GRPC_PORT=8085             \
MONETASA_STREAMS_DB_URL=localhost           \
MONETASA_STREAMS_DB_NAME=monetasa-streams   \
MONETASA_AUTH_URL=localhost:8081            \
$BUILD/monetasa-streams &

###
# Subscriptions
###
MONETASA_SUBSCRIPTIONS_PORT=8086                        \
MONETASA_SUBSCRIPTIONS_DB_URL=localhost                 \
MONETASA_SUBSCRIPTIONS_DB_NAME=monetasa-subscriptions   \
MONETASA_AUTH_URL=localhost:8081                        \
$BUILD/monetasa-subscriptions &

###
# UI
###
cd $UI
npm start &
cd -

trap cleanup EXIT

while : ; do sleep 1 ; done
