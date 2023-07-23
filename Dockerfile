FROM golang:1.19.0

WORKDIR /usr/src/app

RUN go install github.com/cosmtrek/air@latest

COPY go.mod go.sum ./

RUN go mod download 

RUN go mod tidy

EXPOSE 8080
