
FROM golang:1.13-alpine AS builder
ARG SVC_NAME
WORKDIR /go/src/github.com/datapace/datapace
COPY . .
RUN apk update \
    && apk add make \
    && make $SVC_NAME \
    && mv build/datapace-$SVC_NAME /exe

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /etc/nsswitch.conf /etc/nsswitch.conf
COPY --from=builder /exe /
ENTRYPOINT ["/exe"]
