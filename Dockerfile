FROM --platform=$BUILDPLATFORM golang:1.18.2-alpine as builder

ENV CGO_ENABLED=1

# TODO: remove
ARG CI_SERVER_HOST=github.com
ARG CI_REGISTRY_USER=malekvictor
ARG CI_JOB_TOKEN=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


ARG TARGETOS
ARG TARGETARCH

RUN apk update && apk add --no-cache make git build-base musl-dev librdkafka librdkafka-dev
WORKDIR /go/src/github.com/space-box-writer
COPY . ./

RUN echo "machine ${CI_SERVER_HOST} login ${CI_REGISTRY_USER} password ${CI_JOB_TOKEN}" > ~/.netrc

RUN echo "build binary on os: $TARGETOS for platform: $TARGETARCH" && \
    export PATH=$PATH:/usr/local/go/bin && \
    export GOPRIVATE=github.com/hexy-dev/space-box/ && \
    go mod download && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -tags musl /go/src/github.com/space-box-writer/cmd/main.go && \
    mv main /space-box-writer/main && \
    rm -Rf /usr/local/go/src

RUN echo "move migration folder" && \
    mkdir -p /space-box-writer/migrations && \
    mv /go/src/github.com/space-box-writer/adapter/clickhouse/migrations /space-box-writer/migrations


FROM alpine:latest as app
WORKDIR /space-box-writer
COPY --from=builder /space-box-writer/. /space-box-writer/
CMD ./main