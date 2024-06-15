FROM golang:1.20.14

WORKDIR /app

RUN apt-get update && \
    apt-get -y --no-install-recommends install software-properties-common && \
    add-apt-repository "deb http://httpredir.debian.org/debian bullseye main" && \
    apt-get update && \
    apt-get -qq install -y libvips-dev && rm -rf /var/lib/apt/lists/*

COPY components/proxy/ ./

COPY datasets/model.onnx ./model.onnx
COPY datasets/model-weight.onnx ./model-weight.onnx

RUN go mod tidy

RUN go mod download

RUN export CGO_ENABLED=1

RUN go build -o proxy

EXPOSE 8082

ENTRYPOINT [ "./proxy" ]