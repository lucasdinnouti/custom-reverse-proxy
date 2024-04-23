FROM golang:1.19.3-alpine

WORKDIR /app

COPY components/runner/* ./

COPY datasets/chat_dos/chat_dos.txt ./testcase_1.txt

RUN go mod download

RUN go build -o runner

EXPOSE 8081

ENTRYPOINT [ "./runner" ]