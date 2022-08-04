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

| Variable                       | Description                              | Default               |
|--------------------------------|------------------------------------------|-----------------------|
| DATAPACE_SUBSCRIPTIONS_PORT    | Stream service port                      | localhost             |
| DATAPACE_SUBSCRIPTIONS_DB_URL  | List of database cluster URLs            | 0.0.0.0               |
| DATAPACE_SUBSCRIPTIONS_DB_NAME | Name of the database used by the service | subscriptions         |
| DATAPACE_SUBSCRIPTIONS_DB_USER | Database user                            |                       |
| DATAPACE_SUBSCRIPTIONS_DB_PASS | Database password                        |                       |
| DATAPACE_AUTH_URL              | Auth service gRPC URL                    | localhost:8081        |
| DATAPACE_TRANSACTIONS_URL      | Transactions service gRPC URL            | localhost:8081        |
| DATAPACE_STREAMS_URL           | Streams service gRPC URL                 | localhost:8081        |
| DATAPACE_SHARING_URL           | Sharing service gRPC URL (optional)      | localhost:8081        |
| DATAPACE_PROXY_URL             | Proxy service URL                        | http://localhost:8080 |

## Deployment

The service itself is distributed as Docker container. You can find a Docker composition
[here](../docker/docker-compose.yml).

To start the service outside of the container, execute the following shell script:

```bash
cd $GOPATH/src/datapace

# compile the subscriptions
make subscriptions

# copy binary to bin
make install

# set the environment variables and run the service
DATAPACE_SUBSCRIPTIONS_PORT=[Service port] DATAPACE_SUBSCRIPTIONS_DB_URL=[List of database cluster URLs] DATAPACE_SUBSCRIPTIONS_DB_NAME=[Name of the database used by the service] DATAPACE_SUBSCRIPTIONS_DB_USER=[Database user] DATAPACE_SUBSCRIPTIONS_DB_PASS=[Database password] DATAPACE_AUTH_URL=[Auth service gRPC URL] DATAPACE_PROXY_URL=[Proxy service URL] GOOGLE_APPLICATION_CREDENTIALS=[Path to Google app credentials file] $GOBIN/datapace-subscriptions
```

Please note the presence of the GOOGLE_APPLICATION_CREDENTIALS env variable which is used by Google Big Query client.

## Usage

For more information about service capabilities and its usage, please check out
the [API documentation](swagger.yml).
