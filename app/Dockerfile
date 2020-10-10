# Initial stage: download modules
FROM golang:1.15 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

# Intermediate stage: Build the binary
FROM golang:1.15 as builder
COPY --from=modules /go/pkg /go/pkg

# add a non-privileged user
RUN useradd -u 10001 trip

RUN mkdir -p /trip-api
ADD . /trip-api
WORKDIR /trip-api

# Build the binary with go build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o ./bin/trip-api ./

# Final stage: Run the binary
FROM scratch

# don't forget /etc/passwd from previous stage
COPY --from=builder /etc/passwd /etc/passwd
USER trip

# and finally the binary
COPY --from=builder /trip-api/bin/trip-api /trip-api

CMD ["/trip-api"]