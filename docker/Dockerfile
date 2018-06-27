FROM golang:1.10-alpine AS builder
ARG SVC_NAME
WORKDIR /go/src/monetasa
COPY . .
RUN apk update \
    && apk add git gcc libc-dev libltdl libtool libgcc make \
    && make $SVC_NAME \
    && mv build/monetasa-$SVC_NAME /exe

FROM golang:1.10-alpine
RUN apk update && apk add libltdl
COPY --from=builder /exe /
ENTRYPOINT ["/exe"]
