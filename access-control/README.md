# Access Control

Access control service provides an HTTP API for controlling access to streams.
Through this API clients are able request, fetch, approve and revoke access
requests.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                             | Description                              | Default                                    |
| ------------------------------------ | ---------------------------------------- | ------------------------------------------ |
| DATAPACE_ACCESS_CONTROL_HTTP_PORT    | Access control service HTTP port         | 8080                                       |
| DATAPACE_ACCESS_CONTROL_GRPC_PORT    | Access control service gRPC port         | 8081                                       |
| DATAPACE_ACCESS_CONTROL_DB_URL       | Database URL                             | 0.0.0.0                                    |
| DATAPACE_ACCESS_CONTROL_DB_USER      | Database user                            |                                            |
| DATAPACE_ACCESS_CONTROL_DB_PASS      | Database password                        |                                            |
| DATAPACE_ACCESS_CONTROL_DB_NAME      | Name of the database used by the service | access                                     |
| DATAPACE_AUTH_URL                    | Auth service gRPC URL                    | localhost:8081                             |
| DATAPACE_ACCESS_CONTROL_FABRIC_ADMIN | Organization admin for Fabric            | admin                                      |
| DATAPACE_ACCESS_CONTROL_FABRIC_NAME  | Organization name for Fabric             | org1                                       |
| DATAPACE_CONFIG                      | Path to the configuration directory      | `/src/github.com/datapace/datapace/config` |
| DATAPACE_ACCESS_CONTROL_CHAINCODE    | Access Control chaincode ID              | access                                     |

## Deployment

The service itself is distributed as Docker container. You can find a Docker composition
[here](../docker/docker-compose.yml).

To start the service outside of the container, execute the following shell script:

```bash
cd $GOPATH/src/datapace

# compile the access-control
make access-control

# copy binary to bin
make install

# set the environment variables and run the service
DATAPACE_ACCESS_CONTROL_HTTP_PORT=[Access control service HTTP port] DATAPACE_ACCESS_CONTROL_GRPC_PORT=[Access control service gRPC port] DATAPACE_ACCESS_CONTROL_DB_URL=[Database URL] DATAPACE_ACCESS_CONTROL_DB_USER=[Database user] DATAPACE_ACCESS_CONTROL_DB_PASS=[Database password] DATAPACE_ACCESS_CONTROL_DB_NAME=[Name of the database used by the service] DATAPACE_AUTH_URL=[Auth service gRPC URL] DATAPACE_ACCESS_CONTROL_FABRIC_ADMIN=[Organization admin for Fabric] DATAPACE_ACCESS_CONTROL_FABRIC_NAME=[Organization name for Fabric] DATAPACE_CONFIG=[Path to the configuration directory] DATAPACE_ACCESS_CONTROL_CHAINCODE=[Access Control chaincode ID] $GOBIN/datapace-access-control
```

## Usage

For more information about service capabilities and its usage, please check out
the [API documentation](swagger.yml).
