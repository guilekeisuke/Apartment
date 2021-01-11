FROM golang:latest AS build

WORKDIR /go/api
ENV GO111MODULE=on
COPY ./api/go.mod .
COPY ./api/go.sum .
RUN go mod download
COPY ./api .
RUN go get -u github.com/cyruzin/golang-tmdb

# RUN go build -o /go/bin/app && \
#   go get gopkg.in/urfave/cli.v2@v2.3.0 && \
#   go get github.com/oxequa/realize 