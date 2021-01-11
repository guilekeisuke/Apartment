FROM golang:latest AS build

WORKDIR /go/api
ENV GO111MODULE=on
COPY ./api/go.mod .
COPY ./api/go.sum .
RUN go mod download
COPY ./api .

# RUN go build -o /go/bin/app && \
#   go get gopkg.in/urfave/cli.v2@master && \
#   go get github.com/oxequa/realize && \
#   go get -u github.com/cyruzin/golang-tmdb