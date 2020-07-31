# Executions

Executions service provides an HTTP API for managing executions of algorithm
over datasets.
Through this API clients are able to create and fetch executions.

## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                      | Description                              | Default                |
|-------------------------------|------------------------------------------|------------------------|
| DATAPACE_EXECUTIONS_HTTP_PORT | Executions service HTTP port             | 8080                   |
| DATAPACE_EXECUTIONS_DB_URL    | List of database cluster URLs            | 0.0.0.0                |
| DATAPACE_EXECUTIONS_DB_USER   | Database user                            |                        |
| DATAPACE_EXECUTIONS_DB_PASS   | Database password                        |                        |
| DATAPACE_EXECUTIONS_DB_NAME   | Name of the database used by the service | executions             |
| DATAPACE_AUTH_URL             | Auth service gRPC URL                    | localhost:8081         |
| DATAPACE_AI_SYSTEM            | AI system in use (kubeflow,wwf,argo)     | kubeflow               |

## Deployment

The service itself is distributed as Docker container. You can find a Docker composition
[here](../docker/docker-compose.yml).

To start the service outside of the container, execute the following shell script:

```bash
cd $GOPATH/src/datapace

# compile the executions
make executions

# copy binary to bin
make install

# set the environment variables and run the service
DATAPACE_EXECUTIONS_HTTP_PORT=[Executions service HTTP port] DATAPACE_EXECUTIONS_DB_URL=[List of database cluster URLs] DATAPACE_EXECUTIONS_DB_USER=[Database user] DATAPACE_EXECUTIONS_DB_PASS=[Database password] DATAPACE_EXECUTIONS_DB_NAME=[Name of the database used by the service] DATAPACE_AUTH_URL=[Auth service gRPC URL] $GOBIN/datapace-executions
```

## Usage

For more information about service capabilities and its usage, please check out
the [API documentation](swagger.yml).

## Argo workflows preparation

Argo pipeline example: 

At this moment namespace in use is `default` and template name is `datapace`
This is example content of the metadata field which is supplied to datapace when new algo is created.
```

 {
     "pipeline": "{\"template\":{\"metadata\":{\"name\":\"TEMPLATENAMEHERE\",\"namespace\":\"default\"},\"spec\":{\"templates\":[{\"name\":\"datapace\",\"inputs\":{\"parameters\":[{\"name\":\"datapace_url\"}]},\"container\":{\"name\":\"main\",\"image\":\"docker/whalesay\",\"command\":[\"cowsay\"],\"args\":[\"{{inputs.parameters.datapace_url}}\"]}}]}}}"
  }
```
