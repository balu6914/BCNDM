FROM golang:1.8-alpine AS builder
WORKDIR /go/src/monetasa
COPY . .
RUN cd cmd && CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o exe

FROM scratch
COPY --from=builder /go/src/monetasa/cmd/monetasa/exe /
EXPOSE 8080
ENTRYPOINT ["/exe"]
