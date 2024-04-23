FROM golang:1.19.3-alpine

WORKDIR /app

COPY components/processor/* ./

RUN go mod download

RUN go build -o processor

EXPOSE 8083

ENTRYPOINT [ "./processor" ]