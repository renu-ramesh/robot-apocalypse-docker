
# syntax=docker/dockerfile:1
# Alpine is chosen for its small footprint
# compared to Ubuntu
FROM golang:1.16-alpine

WORKDIR /app
# Download necessary Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go get github.com/renu-ramesh/robot-apocalypse-docker/common
RUN go get github.com/renu-ramesh/robot-apocalypse-docker/models
RUN go get github.com/renu-ramesh/robot-apocalypse-docker/mongodb

## we run go build to compile the binary
## executable of our Go program
RUN go build -o main .
## Our start command which kicks off
## our newly created binary executable
CMD ["/app/main"]