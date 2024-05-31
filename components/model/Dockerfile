FROM golang:1.20.14 AS base

WORKDIR /app

RUN apt-get update && \
    apt-get -y --no-install-recommends install software-properties-common && \
    add-apt-repository "deb http://httpredir.debian.org/debian bullseye main" && \
    apt-get update && \
    apt-get -qq install -y libvips-dev && rm -rf /var/lib/apt/lists/*

COPY components/model/ ./

RUN go mod tidy

RUN go mod download

RUN export CGO_ENABLED=1

RUN go build -o model

ENTRYPOINT [ "./model" ]