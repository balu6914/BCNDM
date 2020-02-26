# dProxy

## Description

dProxy serves as a gateway to resources, typically files provided via HTTP or local filesystem.
Common use case for dProxy would be to have it being able to access some protected resources while direct access to those is not possible.
In that use case, dProxy would act as a gatekeeper to a resource.
In order to get the resource from dProxy, requester needs to pass dProxy authorization.
Upon successful authorization, dProxy will fetch the protected resource and deliver it to requester.

## Usage
In order to grant access to the resource behind the dProxy, one needs to create the authorization JWT token.
this is done by doing HTTP POST request to  `/api/register` endpoint with the json payload describing the resource URL and the ttl (number of seconds for which the token is valid).
More information about this endpoint can be found in the swagger.yml file in this directory. 

Once user has the JWT token, he can fetch the proxied resource by accessing `/http` endpoint on dProxy.
User should issue HTTP GET request to `/http` endpoint with HTTP header: "Authorization: your-jwt-token-here".
Upon this request, dProxy will analyze the JWT token from the HTTP header, and if valid, it will fetch the resource and send it to user.


## Configuration


The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                             | Description                              | Default                |
|--------------------------------------|------------------------------------------|------------------------|
| DATAPACE_PROXY_HTTP_PORT             | Reverse proxy HTTP port                  | 9090                   |
| DATAPACE_JWT_SECRET                  | Reverse proxy JWT secret                 | examplesecret       |

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
DATAPACE_PROXY_HTTP_PORT=[Reverse proxy HTTP port] DATAPACE_JWT_SECRET=[Reverse proxy JWT secret] $GOBIN/datapace-dproxy
```