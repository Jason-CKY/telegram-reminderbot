FROM node:18-alpine as dev

WORKDIR /app

RUN wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz &&\
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz &&\
    rm go1.22.0.linux-amd64.tar.gz

ENV PATH="${PATH}:/usr/local/go/bin:/root/go/bin"

RUN go install github.com/cosmtrek/air@latest

ENTRYPOINT ["air"]

FROM golang:1.21-alpine as production

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /main
CMD [ "/main" ]
