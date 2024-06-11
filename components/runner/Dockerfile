FROM golang:1.20.14-alpine

WORKDIR /app

COPY components/runner/ ./

COPY datasets/chat_dos/shuffled.txt ./testcase.txt

RUN go mod tidy

RUN go mod download

RUN go build -o runner

EXPOSE 8081

ENTRYPOINT [ "./runner" ]