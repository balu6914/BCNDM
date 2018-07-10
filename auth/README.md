# Auth

Auth service provides an HTTP API for authorization and managing users.
Through this API clients are able to get and update user info, register and
login.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                  | Description                              | Default        |
|---------------------------|------------------------------------------|----------------|
| MONETASA_AUTH_HTTP_PORT   | Auth service HTTP port                   | 8080           |
| MONETASA_AUTH_GRPC_PORT   | Auth service gRPC port                   | 8081           |
| MONETASA_AUTH_MONGO_URL   | List of database cluster URLs            | 0.0.0.0        |
| MONETASA_AUTH_MONGO_USER  | Database user                            |                |
| MONETASA_AUTH_MONGO_PASS  | Database password                        |                |
| MONETASA_AUTH_MONGO_DB    | Name of the database used by the service | auth           |
| MONETASA_TRANSACTIONS_URL | Transactions service gRPC URL            | localhost:8081 |
| MONETASA_AUTH_SECRET      | Authorization secret                     | monetasa       |

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
MONETASA_AUTH_HTTP_PORT=[Auth service HTTP port] MONETASA_AUTH_GRPC_PORT=[Auth service gRPC port] MONETASA_AUTH_MONGO_URL=[List of database cluster URLs] MONETASA_AUTH_MONGO_USER=[Database user] MONETASA_AUTH_MONGO_PASS=[Database password] MONETASA_AUTH_MONGO_DB=[Name of the database used by the service] MONETASA_TRANSACTIONS_URL=[Transactions service gRPC URL] MONETASA_AUTH_SECRET=[Authorization secret] $GOBIN/monetasa-auth
```

## Usage

For more information about service capabilities and its usage, please check out
the [API documentation](swagger.yml).
