FROM golang:1-alpine as builder

RUN apk --no-cache --no-progress add git ca-certificates tzdata make \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

WORKDIR /go/echoserver

# Download go modules
COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://proxy.golang.org go mod download

COPY . .

RUN make build

# Create a minimal container to run a Golang static binary
FROM scratch

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/echoserver .

ENV ECHO_HTTP1_PORT=8001
ENV ECHO_HTTP2_CLEARTEXT_PORT=8002

ENTRYPOINT ["/echoserver", "-verbose"]
EXPOSE 8001/tcp
EXPOSE 8002/tcp
EXPOSE 8003/tcp