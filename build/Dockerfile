FROM golang:1.22-alpine as dev

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

ENTRYPOINT ["air"]

FROM golang:1.22-alpine as production

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /main
CMD [ "/main" ]
