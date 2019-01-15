# Transactions

Transactions service provides an HTTP API for managing transactions and balance.
Through this API clients are able to get users balance and initiate
transactions.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                               | Description                              | Default                |
|----------------------------------------|------------------------------------------|------------------------|
| DATAPACE_TRANSACTIONS_HTTP_PORT        | Transactions service HTTP port           | 8080                   |
| DATAPACE_TRANSACTIONS_GRPC_PORT        | Transactions service gRPC port           | 8081                   |
| DATAPACE_TRANSACTIONS_DB_URL           | List of database cluster URLs            | 0.0.0.0                |
| DATAPACE_TRANSACTIONS_DB_USER          | Database user                            |                        |
| DATAPACE_TRANSACTIONS_DB_PASS          | Database password                        |                        |
| DATAPACE_TRANSACTIONS_DB_NAME          | Name of the database used by the service | transactions           |
| DATAPACE_TRANSACTIONS_FABRIC_ADMIN     | Organization admin for Fabric            | admin                  |
| DATAPACE_TRANSACTIONS_FABRIC_NAME      | Organization name for Fabric             | org1                   |
| DATAPACE_CONFIG                        | Path to the configuration directory      | `/src/datapace/config` |
| DATAPACE_TRANSACTIONS_FABRIC_CHAINCODE | Fabric token chaincode id                | token                  |
| DATAPACE_AUTH_URL                      | Auth service gRPC URL                    | localhost:8081         |

## Deployment

The service itself is distributed as Docker container. You can find a Docker composition
[here](../docker/docker-compose.yml).

To start the service outside of the container, execute the following shell script:

```bash
cd $GOPATH/src/datapace

# compile the transactions
make transactions

# copy binary to bin
make install

# set the environment variables and run the service
DATAPACE_TRANSACTIONS_HTTP_PORT=[Transactions service HTTP port] DATAPACE_TRANSACTIONS_GRPC_PORT=[Transactions service gRPC port] DATAPACE_TRANSACTIONS_DB_URL=[List of database cluster URLs] DATAPACE_TRANSACTIONS_DB_USER=[Database user] DATAPACE_TRANSACTIONS_DB_PASS=[Database password] DATAPACE_TRANSACTIONS_DB_NAME=[Name of the database used by the service] DATAPACE_TRANSACTIONS_FABRIC_ADMIN=[Organization admin for Fabric] DATAPACE_TRANSACTIONS_FABRIC_NAME=[Organization name for Fabric] DATAPACE_CONFIG=[Fabric configuration directory path] DATAPACE_TRANSACTIONS_FABRIC_CHAINCODE=[Fabric token chaincode id] DATAPACE_AUTH_URL=[Auth service gRPC URL] $GOBIN/datapace-transactions
```

**IMPORTANT:** _Please note that DATAPACE_CONFIG env variable is also used in [Fabric config](../config/fabric/config.yaml) as a path to the configuration directory. Path should be provided without trailing `/` character._

## Usage

For more information about service capabilities and its usage, please check out
the [API documentation](swagger.yml).
