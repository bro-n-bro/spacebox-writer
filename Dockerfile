FROM --platform=$BUILDPLATFORM golang:1.18.2-alpine as builder

ENV CGO_ENABLED=1

ARG TARGETOS
ARG TARGETARCH

RUN apk update && apk add --no-cache make git build-base musl-dev librdkafka librdkafka-dev
WORKDIR /go/src/github.com/space-box-writer
COPY . ./

RUN echo "build binary on os: $TARGETOS for platform: $TARGETARCH" && \
    export PATH=$PATH:/usr/local/go/bin && \
    go mod download && \
    GOOS=$TARGETOS GOARCH=$TARGETARCH go build -tags musl /go/src/github.com/spacebox-writer/cmd/main.go && \
    mv main /spacebox-writer/main && \
    rm -Rf /usr/local/go/src

RUN echo "move migration folder" && \
    mkdir -p /spacebox-writer/migrations && \
    mv /go/src/github.com/spacebox-writer/adapter/clickhouse/migrations /spacebox-writer/migrations


FROM alpine:latest as app
WORKDIR /spacebox-writer
COPY --from=builder /spacebox-writer/. /spacebox-writer/
CMD ./main
