# Terms

Terms service provides an HTTP and gRPC API for managing terms of use for streams.
Through this API clients are able to do the following
actions:

- create terms hash for a stream and store it in blockchain (gRPC)
- validate terms hash for a stream (HTTP)

## Configuration

The service is configured using the environment variables presented in the
following table. Note that for any missing values, the defaults will be used.

| Variable                    | Description                              | Default                                  |
|-----------------------------|------------------------------------------|------------------------------------------|
| DATAPACE_TERMS_HTTP_PORT    | Terms service HTTP port                  | 8080                                     |
| DATAPACE_TERMS_GRPC_PORT    | Terms service gRPC port                  | 8081                                     |
| DATAPACE_TERMS_DB_URL       | List of database cluster URLs            | 0.0.0.0                                  |
| DATAPACE_TERMS_DB_NAME      | Name of the database used by the service | subscriptions                            |
| DATAPACE_TERMS_DB_USER      | Database user                            |                                          |
| DATAPACE_TERMS_DB_PASS      | Database password                        |                                          |
| DATAPACE_AUTH_URL           | Auth service gRPC URL                    | localhost:8081                           |
| DATAPACE_TERMS_FABRIC_ADMIN | Fabric admin user                        | admin                                    |
| DATAPACE_TERMS_FABRIC_NAME  | Fabric organization                      | org1                                     |
| DATAPACE_CONFIG             | Path to the datapace configuration       | /src/github.com/datapace/datapace/config |
| DATAPACE_TERMS_CHAINCODE    | Name of the terms chaincode in Fabric    | terms                                    |

## Usage

When the CreateTerms function is invoked via gRPC (usually by the stream service), the term service will try to fetch
the terms of use document for the stream, from the supplied location, it will then create document hash and store that
hash in it's database and blockchain.

When validate is invoked via HTTP (GET request to /streams/{StreamID}/terms/{TermsHash}), the service will use the stream id
and the terms hash (from parameters supplied) and it will check the hash against what's written in the blockchain in order to validate it.

For more information about service capabilities and its usage, please see the [API documentation](swagger.yaml).
