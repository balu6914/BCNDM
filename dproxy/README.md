# dProxy

## Configuration


The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                             | Description                              | Default                |
|--------------------------------------|------------------------------------------|------------------------|
| DATAPACE_PROXY_HTTP_PORT             | Reverse proxy HTTP port                  | 9090                   |
| DATAPACE_PROXY_TARGET_URL            | Reverse proxy target URL                 | http://localhost       |

The service itself is distributed as Docker container. You can find a Docker composition
[here](../docker/docker-compose.yml).

To start the service outside of the container, execute the following shell script:

```bash
cd $GOPATH/src/datapace

# compile the dproxy
make dproxy

# copy binary to bin
make install

# set the environment variables and run the service
DATAPACE_PROXY_HTTP_PORT=[Reverse proxy HTTP port] DATAPACE_PROXY_TARGET_URL=[Reverse proxy target URL] $GOBIN/datapace-dproxy
```