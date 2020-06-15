# Auth

Auth service provides an HTTP API for authorization and managing users.
Through this API clients are able to get and update user info, register and
login.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                    | Description                                  | Default                  |
|-----------------------------|----------------------------------------------|--------------------------|
| DATAPACE_AUTH_HTTP_PORT     | Auth service HTTP port                       | 8080                     |
| DATAPACE_AUTH_GRPC_PORT     | Auth service gRPC port                       | 8081                     |
| DATAPACE_AUTH_DB_URL        | List of database cluster URLs                | 0.0.0.0                  |
| DATAPACE_AUTH_DB_USER       | Database user                                |                          |
| DATAPACE_AUTH_DB_PASS       | Database password                            |                          |
| DATAPACE_AUTH_DB_NAME       | Name of the database used by the service     | auth                     |
| DATAPACE_TRANSACTIONS_URL   | Transactions service gRPC URL                | localhost:8081           |
| DATAPACE_ACCESS_CONTROL_URL | Access control service gRPC URL              | localhost:8081           |
| DATAPACE_AUTH_SECRET        | Authorization secret                         | datapace                 |
| DATAPACE_ADMIN_EMAIL        | Email of the initial admin for the service   | admin@datapace.localhost |
| DATAPACE_ADMIN_PASSWORD     | Password of the initial admin for the service| datapaceadmin            |

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
DATAPACE_AUTH_HTTP_PORT=[Auth service HTTP port] DATAPACE_AUTH_GRPC_PORT=[Auth service gRPC port] DATAPACE_AUTH_DB_URL=[List of database cluster URLs] DATAPACE_AUTH_DB_USER=[Database user] DATAPACE_AUTH_DB_PASS=[Database password] DATAPACE_AUTH_DB_NAME=[Name of the database used by the service] DATAPACE_TRANSACTIONS_URL=[Transactions service gRPC URL] DATAPACE_ACCESS_CONTROL_URL=[Access control service gRPC URL] DATAPACE_AUTH_SECRET=[Authorization secret] $GOBIN/datapace-auth
```

During the service start, if not present, datapace admin will be created.

## Usage

For more information about service capabilities and its usage, please check out
the [API documentation](swagger.yml).
