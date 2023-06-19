
build:
	CGO_ENABLED=0 go build -a --trimpath --installsuffix cgo --ldflags="-s" -o echoserver

test:
	go test -v -cover ./...

check:
	golangci-lint run

image: check test
	docker build --tag hattorious/echoserver .

run: image
	docker run --rm -p 8001-8005:8001-8005 hattorious/echoserver

.PHONY: check test build image run