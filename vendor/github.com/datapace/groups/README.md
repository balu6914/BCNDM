# Groups

Groups service provides HTTP and gRPC API for groups and group-specific user operations

## 1. Build

```shell
cd $GOPATH/src/groups
make build
```

## 2. Configuration

| Variable                  | Description                              | Default        |
|---------------------------|------------------------------------------|----------------|
| DATAPACE_GROUPS_HTTP_PORT | Service HTTP port                        | 8080           |
| DATAPACE_GROUPS_GRPC_PORT | Service gRPC port                        | 8081           |
| DATAPACE_GROUPS_DB_URL    | Database URL                             | 0.0.0.0        |
| DATAPACE_GROUPS_DB_USER   | Database user                            |                |
| DATAPACE_GROUPS_DB_PASS   | Database password                        |                |
| DATAPACE_GROUPS_DB_NAME   | Name of the database used by the service | groups         |
| DATAPACE_AUTH_URL         | Auth service gRPC URL                    | localhost:8081 |
| DATAPACE_SHARING_URL      | Sharing service gRPC URL                 | localhost:8082 |

## 3. Deployment

### 3.1. Local

#### 3.1.1. Standalone

```shell
# set the environment variables and run the service
DATAPACE_GROUPS_HTTP_PORT=8094 \
DATAPACE_GROUPS_GRPC_PORT=8095 \
DATAPACE_GROUPS_DB_URL=localhost:27027 \
DATAPACE_GROUPS_DB_USER= \
DATAPACE_GROUPS_DB_PASS= \
DATAPACE_GROUPS_DB_NAME=groups \
DATAPACE_AUTH_URL=localhost:8081 \
DATAPACE_SHARING_URL=localhost:8082 \
./datapace-groups
```

#### 3.1.2. Docker

1. [Run Datapace using docker-compose](https://github.com/datapace/datapace/blob/master/docker/README.md)
2. [Run Datapace Sharing service](https://github.com/datapace/sharing/tree/init#3-deployment)
3. Build the groups docker image:
    ```shell 
    make docker
    ```
4. Run the groups DB:
    ```shell
    make rundb
    ```
5. Run groups service:
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
2. gRPC: see [proto/groups.proto](proto/groups.proto) file for the reference
