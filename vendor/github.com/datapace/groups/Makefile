.PHONY: test proto clean
default: build

BINARY_FILE_NAME=datapace-groups
COVERAGE_FILE_NAME=cover.out
GOPRIVATE=github.com/datapace/*
THRESHOLD=85

vet:
	go vet

test: vet
	go test -race -cover -coverprofile=${COVERAGE_FILE_NAME} ./...
	cat ${COVERAGE_FILE_NAME} | grep -v mock > ${COVERAGE_FILE_NAME}.tmp
	mv -f ${COVERAGE_FILE_NAME}.tmp ${COVERAGE_FILE_NAME}
	go tool cover -func=${COVERAGE_FILE_NAME} | grep -Po '^total\:\h+\(statements\)\h+\K[\d]+' > cover.tmp
	./cover.sh

proto:
	go get google.golang.org/grpc@v1.27.0
	go install github.com/golang/protobuf/protoc-gen-go@v1.3.4
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@38aafd89f814f347db56a52efd055961651078ad
	PATH=~/go/bin:$(PATH) protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative proto/*.proto

build:
	git config --global url."git@github.com:datapace/datapace".insteadOf "https://github.com/datapace/datapace"
	git config --global url."git@github.com:datapace/sharing".insteadOf "https://github.com/datapace/sharing"
	CGO_ENABLED=0 GOOS=linux GOARCH= GOARM= go build -o ${BINARY_FILE_NAME} cmd/main.go
	chmod ugo+x ${BINARY_FILE_NAME}

docker: build
	docker build -t datapace/groups .

rundb:
	docker run \
		-d \
		--name datapace-groups-db \
		--network docker_datapace-net \
		-p 27027:27017 \
		-v data:/data/db \
		mongo:latest \
		--nojournal

run:
	docker run \
		-d \
		--name datapace-groups \
		--network docker_datapace-net \
		-p 8094:8094 \
		--expose 8095 \
		--env DATAPACE_GROUPS_HTTP_PORT=8094 \
		--env DATAPACE_GROUPS_GRPC_PORT=8095 \
		--env DATAPACE_GROUPS_DB_URL=datapace-groups-db \
		--env DATAPACE_AUTH_URL=datapace-auth:8081 \
		--env DATAPACE_SHARING_URL=datapace-sharing:8097 \
		datapace/groups

clean:
	go clean
	rm -f ${BINARY_FILE_NAME} ${COVERAGE_FILE_NAME}
