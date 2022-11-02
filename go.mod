module github.com/datapace/datapace

go 1.14

require (
	cloud.google.com/go v0.34.0
	github.com/asaskevich/govalidator v0.0.0-20180315120708-ccb8e960c48f
	github.com/datapace/events v0.1.9
	github.com/datapace/groups v0.1.0
	github.com/datapace/sharing v0.1.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-kit/kit v0.9.0
	github.com/go-zoo/bone v1.3.0
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/websocket v1.4.0
	github.com/hyperledger/fabric v1.4.0
	github.com/hyperledger/fabric-sdk-go v1.0.0-beta2
	github.com/jmoiron/sqlx v1.2.0
	github.com/johnfercher/maroto v0.27.0
	github.com/lib/pq v1.3.0
	github.com/oklog/ulid/v2 v2.0.2
	github.com/ory/dockertest v3.3.4+incompatible
	github.com/prometheus/client_golang v1.11.1
	github.com/rubenv/sql-migrate v0.0.0-20200212082348-64f95ea68aa3
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.7.1
	golang.org/x/crypto v0.0.0-20220315160706-3147a52a75dd
	google.golang.org/api v0.1.0
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
)

replace (
	github.com/go-kit/kit => github.com/go-kit/kit v0.8.0
	github.com/go-zoo/bone => github.com/go-zoo/bone v0.0.0-20190117145001-d7ce1372afa7
	github.com/gogo/protobuf => github.com/gogo/protobuf v1.2.0
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
	golang.org/x/crypto => golang.org/x/crypto v0.0.0-20200204104054-c9f3fb736b72
)
