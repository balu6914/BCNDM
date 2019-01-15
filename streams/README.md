# Streams

Streams service provides an HTTP API for managing streams.
Through this API clients are able to do the following
actions:

- create, update, retrieve and delete a stream
- search streams
- create new bulk of streams

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                   | Description                              | Default        |
|----------------------------|------------------------------------------|----------------|
| DATAPACE_STREAMS_HTTP_PORT | Stream service HTTP port                 | 8080           |
| DATAPACE_STREAMS_GRPC_PORT | Stream service gRPC port                 | 8081           |
| DATAPACE_STREAMS_DB_URL    | List of database cluster URLs            | 0.0.0.0        |
| DATAPACE_STREAMS_DB_NAME   | Name of the database used by the service | streams        |
| DATAPACE_STREAMS_DB_USER   | Database user                            |                |
| DATAPACE_STREAMS_DB_PASS   | Database password                        |                |
| DATAPACE_AUTH_URL          | Auth service gRPC URL                    | localhost:8081 |

## Deployment

The service itself is distributed as Docker container. You can find a Docker composition
[here](../docker/docker-compose.yml).

To start the service outside of the container, execute the following shell script:

```bash
cd $GOPATH/src/datapace

# compile the streams
make streams

# copy binary to bin
make install

# set the environment variables and run the service
DATAPACE_STREAMS_HTTP_PORT=[Service HTTP port] DATAPACE_STREAMS_GRPC_PORT=[Service gRPC port] DATAPACE_STREAMS_DB_URL=[List of database cluster URLs] DATAPACE_STREAMS_DB_NAME=[Name of the database used by the service] DATAPACE_STREAMS_DB_USER=[Database user] DATAPACE_STREAMS_DB_PASS=[Database password] DATAPACE_AUTH_URL=[Auth service gRPC URL] GOOGLE_APPLICATION_CREDENTIALS=[Path to Google app credentials file] $GOBIN/datapace-streams
```

Please note the presence of the GOOGLE_APPLICATION_CREDENTIALS env variable which is used by Google Big Query client.

## Usage

For more information about service capabilities and its usage, please check out
the [API documentation](swagger.yml).
