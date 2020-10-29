# Initial stage: download modules
FROM golang:1.15 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

# Intermediate stage: Build the binary
FROM golang:1.15 as builder
COPY --from=modules /go/pkg /go/pkg

RUN mkdir -p /trip-api
ADD . /trip-api
WORKDIR /trip-api

# Build the binary with go build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o ./bin/trip-api ./

CMD ["/trip-api/bin/trip-api"]