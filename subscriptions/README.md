# Subscriptions

Subscriptions service provides an HTTP API for managing subscriptions.
Through this API clients are able to do the following
actions:

- create, update, retrieve and delete a subscription
- search subscriptions

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                       | Description                              | Default        |
|--------------------------------|------------------------------------------|----------------|
| MONETASA_SUBSCRIPTIONS_PORT    | Stream service port                      | localhost      |
| MONETASA_SUBSCRIPTIONS_DB_URL  | List of database cluster URLs            | 0.0.0.0        |
| MONETASA_SUBSCRIPTIONS_DB_NAME | Name of the database used by the service | subscriptions  |
| MONETASA_SUBSCRIPTIONS_DB_USER | Database user                            |                |
| MONETASA_SUBSCRIPTIONS_DB_PASS | Database password                        |                |
| MONETASA_AUTH_URL              | Auth service gRPC URL                    | localhost:8081 |
| MONETASA_TRANSACTIONS_URL      | Transactions service gRPC URL            | localhost:8081 |
| MONETASA_STREAMS_URL           | Streams service gRPC URL                 | localhost:8081 |

## Deployment

The service itself is distributed as Docker container. You can find a Docker composition
[here](../docker/docker-compose.yml).

To start the service outside of the container, execute the following shell script:

```bash
cd $GOPATH/src/monetasa

# compile the subscriptions
make subscriptions

# copy binary to bin
make install

# set the environment variables and run the service
MONETASA_SUBSCRIPTIONS_PORT=[Service port] MONETASA_SUBSCRIPTIONS_DB_URL=[List of database cluster URLs] MONETASA_SUBSCRIPTIONS_DB_NAME=[Name of the database used by the service] MONETASA_SUBSCRIPTIONS_DB_USER=[Database user] MONETASA_SUBSCRIPTIONS_DB_PASS=[Database password] MONETASA_AUTH_URL=[Auth service gRPC URL] $GOBIN/monetasa-subscriptions
```

## Usage

For more information about service capabilities and its usage, please check out
the [API documentation](swagger.yml).
