all: build

test:
	golint ./...
	go tool vet *.go
	go test -cover -v ./...

deps:
	go get github.com/golang/lint/golint

build: deps test
	! gofmt -l *.go | read nothing
