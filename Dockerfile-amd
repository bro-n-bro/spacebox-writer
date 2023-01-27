FROM golang:1.18.2-alpine as builder

ENV CGO_ENABLED=1

RUN apk update && apk add --no-cache make git build-base musl-dev librdkafka librdkafka-dev
WORKDIR /go/src/github.com/spacebox-writer
COPY . ./

RUN echo "build binary" && \
    export PATH=$PATH:/usr/local/go/bin && \
    go mod download && \
    go build -tags musl /go/src/github.com/spacebox-writer/cmd/main.go && \
    mkdir -p /spacebox-writer && \
    mv main /spacebox-writer/main && \
    rm -Rf /usr/local/go/src

RUN echo "move migration folder" && \
    mkdir -p /spacebox-writer/migrations && \
    mv /go/src/github.com/spacebox-writer/adapter/clickhouse/migrations/* /spacebox-writer/migrations


FROM alpine:latest as app
WORKDIR /spacebox-writer
COPY --from=builder /spacebox-writer/. /spacebox-writer/
COPY --from=builder /spacebox-writer/migrations/. /spacebox-writer/migrations/
CMD ./main
