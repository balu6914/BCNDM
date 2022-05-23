# Sharing

Sharing service provides HTTP and gRPC API for sharing streams with users and user groups

## 1. Build

```shell
cd $GOPATH/src/sharing
make build
```

## 2. Configuration

| Variable                    | Description                              | Default        |
|-----------------------------|------------------------------------------|----------------|
| DATAPACE_SHARING_HTTP_PORT  | Service HTTP port                        | 8080           |
| DATAPACE_SHARING_GRPC_PORT  | Service gRPC port                        | 8081           |
| DATAPACE_SHARING_DB_URL     | Database URL                             | 0.0.0.0        |
| DATAPACE_SHARING_DB_USER    | Database user                            |                |
| DATAPACE_SHARING_DB_PASS    | Database password                        |                |
| DATAPACE_SHARING_DB_NAME    | Name of the database used by the service | sharing        |
| DATAPACE_AUTH_URL           | Auth service gRPC URL                    | localhost:8081 |

## 3. Deployment

### 3.1. Local

#### 3.1.1. Standalone

```shell
# set the environment variables and run the service
DATAPACE_SHARING_HTTP_PORT=8096 \
DATAPACE_SHARING_GRPC_PORT=8097 \
DATAPACE_SHARING_DB_URL=localhost:27028 \
DATAPACE_SHARING_DB_USER= \
DATAPACE_SHARING_DB_PASS= \
DATAPACE_SHARING_DB_NAME=sharing \
DATAPACE_AUTH_URL=localhost:8081 \
./datapace-sharing
```

#### 3.1.2. Docker

1. [Run Datapace using docker-compose](https://github.com/datapace/datapace/blob/master/docker/README.md)
2. Optionally [run Datapace Groups](https://github.com/datapace/groups/blob/master/README.md#312-docker)
3. Build the sharing docker image:
    ```shell 
    make docker
    ```
4. Run the sharing DB:
    ```shell
    make rundb
    ```
5. Run sharing service
    ```shell
    make run
    ```

#### 3.1.3. K8s

TODO

### 3.2. Cloud

TODO

## 4. Usage

### 4.1. API

1. HTTP: see [swagger.yaml](swagger.yaml) file for the reference
2. gRPC: see [proto/sharing.proto](proto/sharing.proto) file for the reference
