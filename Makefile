
build:
	CGO_ENABLED=0 go build -a --trimpath --installsuffix cgo --ldflags="-s" -o echoserver

test:
	go test -v -cover ./...

check:
	golangci-lint run

image: check test
	docker build --tag hattorious/echoserver .

run: image
	docker run --rm -p 8001-8006:8001-8006 -p 8007:8007/udp hattorious/echoserver

.PHONY: check test build image run