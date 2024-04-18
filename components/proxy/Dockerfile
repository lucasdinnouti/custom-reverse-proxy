FROM golang:1.19.3-alpine

WORKDIR /app

COPY * ./

RUN go mod download

RUN go build -o proxy

EXPOSE 8082

ENTRYPOINT [ "./proxy" ]