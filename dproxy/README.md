# dProxy

## Description

dProxy serves as a gateway to resources, typically files provided via HTTP or local filesystem.
Common use case for dProxy would be to have it being able to access some protected resources while direct access to those is not possible.
In that use case, dProxy would act as a gatekeeper to a resource.
In order to get the resource from dProxy, requester needs to pass dProxy authorization.
Upon successful authorization, dProxy will fetch the protected resource and deliver it to requester.

## Usage
In order to grant access to the resource behind the dProxy, one needs to create the authorization JWT token.
this is done by doing HTTP POST request to  `/api/token` endpoint with the json payload describing the resource URL and the ttl (number of hours for which the token is valid).
More information about this endpoint can be found in the swagger.yml file in this directory. 

Once user has the JWT token, he can fetch the proxied resource by accessing `/http` endpoint on dProxy.
User should issue HTTP GET request to `/http` endpoint with HTTP header: "Authorization: your-jwt-token-here".
Upon this request, dProxy will analyze the JWT token from the HTTP header, and if valid, it will fetch the resource specified in the token and send it to user.

User can also fetch files from local filesystem directory which is configured using `DATAPACE_LOCAL_FS_ROOT ` env variable.
Files within that directory will be available if user sends HTTP GET request to `/fs` endpoint with HTTP header: "Authorization: your-jwt-token-here".
Upon this request, dProxy will analyze the JWT token from the HTTP header, and if valid, it will fetch local file specified in the token and send it to user.

User can also upload files to local filesystem directory which is configured using `DATAPACE_LOCAL_FS_ROOT ` env variable.
Files will be uploaded if user sends HTTP PUT request to `/fs` endpoint with HTTP header: "Authorization: your-jwt-token-here".
Upon this request, dProxy will analyze the JWT token from the HTTP header, and if valid, it will store HTTP body content to the file specified in the token.

For backwards compatibility with old proxy there is also `/api/register` endpoint which will return URL with the token contained in the url. Users can then use that url to access the resource.

## Configuration


The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                             | Description                                                                                    | Default                |
|--------------------------------------|------------------------------------------------------------------------------------------------|------------------------|
| DATAPACE_PROXY_HTTP_PORT             | Reverse proxy HTTP port                                                                        | 9090                   |
| DATAPACE_JWT_SECRET                  | Reverse proxy JWT secret                                                                       | examplesecret          |
| DATAPACE_LOCAL_FS_ROOT               | Local filesystem directory which serves as root directory when serving local files             | /tmp/test              |
| DATAPACE_DPROXY_DB_HOST              | dProxy database host                                                                           | localhost              |
| DATAPACE_DPROXY_DB_PORT              | dProxy database port                                                                           | 5432                   |
| DATAPACE_DPROXY_DB_USER              | dProxy database username                                                                       | dproxy                 |
| DATAPACE_DPROXY_DB_PASS              | dProxy database password                                                                       | dproxy                 |
| DATAPACE_DPROXY_DB                   | dProxy database name                                                                           | dproxy                 |
| DATAPACE_DPROXY_DB_SSL_MODE          | dProxy database ssl switch                                                                     | disable                |
| DATAPACE_DPROXY_DB_SSL_CERT          | dProxy database certificate                                                                    |                        |
| DATAPACE_DPROXY_DB_SSL_KEY           | dProxy database private key                                                                    |                        |
| DATAPACE_DPROXY_DB_SSL_ROOT_CERT     | dProxy database root certificate                                                               |                        |
| DATAPACE_DPROXY_FS_PATH_PREFIX       | the prefix in the generated URL by dproxy that indicates dproxy should fetch a file            | /fs                    |
| DATAPACE_DPROXY_HTTP_PATH_PREFIX     | the prefix in the generated URL by dproxy that indicates dproxy should fetch a http resource   | /http                  |
| DATAPACE_PROXY_HTTP_PROTO            | protocol that will be put in generated URL                                                     | http                   |
| DATAPACE_PROXY_HTTP_HOST             | host that will be put in the generated URL                                                     | localhost              |



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
DATAPACE_PROXY_HTTP_PORT=[Reverse proxy HTTP port] DATAPACE_JWT_SECRET=[Reverse proxy JWT secret] DATAPACE_LOCAL_FS_ROOT=[Path to local files directory] $GOBIN/datapace-dproxy
```