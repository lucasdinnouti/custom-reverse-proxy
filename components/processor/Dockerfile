FROM golang:1.20.14-alpine

WORKDIR /app

COPY components/processor/* ./

RUN go mod tidy

RUN go mod download

RUN go build -o processor

EXPOSE 8083

ENTRYPOINT [ "./processor" ]