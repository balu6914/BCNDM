# Transactions

Transactions service provides an HTTP API for managing transactions and balance.
Through this API clients are able to get users balance and initiate
transactions.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                               | Description                              | Default                                   |
|----------------------------------------|------------------------------------------|-------------------------------------------|
| MONETASA_TRANSACTIONS_HTTP_PORT        | Transactions service HTTP port           | 8080                                      |
| MONETASA_TRANSACTIONS_GRPC_PORT        | Transactions service gRPC port           | 8081                                      |
| MONETASA_TRANSACTIONS_MONGO_URL        | List of database cluster URLs            | 0.0.0.0                                   |
| MONETASA_TRANSACTIONS_MONGO_USER       | Database user                            |                                           |
| MONETASA_TRANSACTIONS_MONGO_PASS       | Database password                        |                                           |
| MONETASA_TRANSACTIONS_MONGO_DB         | Name of the database used by the service | transactions                              |
| MONETASA_TRANSACTIONS_FABRIC_ADMIN     | Organization admin for Fabric            | admin                                     |
| MONETASA_TRANSACTIONS_FABRIC_NAME      | Organization name for Fabric             | org1                                      |
| MONETASA_TRANSACTIONS_FABRIC_CONF      | Fabric configuration path                | `/src/monetasa/config/fabric/config.yaml` |
| MONETASA_TRANSACTIONS_FABRIC_CHAINCODE | Fabric token chaincode id                | token                                     |
| MONETASA_AUTH_URL                      | Auth service gRPC URL                    | localhost:8081                            |

## Deployment

The service itself is distributed as Docker container. You can find a Docker composition
[here](../docker/docker-compose.yml).

To start the service outside of the container, execute the following shell script:

```bash
cd $GOPATH/src/monetasa

# compile the transactions
make transactions

# copy binary to bin
make install

# set the environment variables and run the service
MONETASA_TRANSACTIONS_HTTP_PORT=[Transactions service HTTP port] MONETASA_TRANSACTIONS_GRPC_PORT=[Transactions service gRPC port] MONETASA_TRANSACTIONS_MONGO_URL=[List of database cluster URLs] MONETASA_TRANSACTIONS_MONGO_USER=[Database user] MONETASA_TRANSACTIONS_MONGO_PASS=[Database password] MONETASA_TRANSACTIONS_MONGO_DB=[Name of the database used by the service] MONETASA_TRANSACTIONS_FABRIC_ADMIN=[Organization admin for Fabric] MONETASA_TRANSACTIONS_FABRIC_NAME=[Organization name for Fabric] MONETASA_TRANSACTIONS_FABRIC_CONF=[Fabric configuration path] MONETASA_TRANSACTIONS_FABRIC_CHAINCODE=[Fabric token chaincode id] MONETASA_AUTH_URL=[Auth service gRPC URL] $GOBIN/monetasa-transactions
```

## Usage

For more information about service capabilities and its usage, please check out
the [API documentation](swagger.yml).
