FROM golang:1.20.14-alpine

WORKDIR /app

COPY components/proxy/ ./

RUN ls

RUN go mod tidy

RUN go mod download

RUN go build -o proxy

EXPOSE 8082

ENTRYPOINT [ "./proxy" ]