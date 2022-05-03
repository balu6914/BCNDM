.PHONY: test proto clean
default: build

BINARY_FILE_NAME=datapace-sharing
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
	CGO_ENABLED=0 GOOS=linux GOARCH= GOARM= go build -o ${BINARY_FILE_NAME} cmd/main.go
	chmod ugo+x ${BINARY_FILE_NAME}

docker: build
	docker build -t datapace/sharing .

rundb:
	docker run \
		-d \
		--name datapace-sharing-db \
		--network docker_datapace-net \
		-p 27028:27017 \
		-v datapace-sharing-data:/data/db \
		mongo:latest \
		--nojournal

run:
	docker run \
		-d \
		--name datapace-sharing \
		--network docker_datapace-net \
		-p 8096:8096 \
		--expose 8097 \
		--env DATAPACE_SHARING_HTTP_PORT=8096 \
		--env DATAPACE_SHARING_GRPC_PORT=8097 \
		--env DATAPACE_SHARING_DB_URL=datapace-sharing-db \
		--env DATAPACE_AUTH_URL=datapace-auth:8081
		datapace/sharing

clean:
	go clean
	rm -f ${BINARY_FILE_NAME} ${COVERAGE_FILE_NAME}
