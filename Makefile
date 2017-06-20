ci: fmt lint test

fmt:
	! gofmt -l *.go | read nothing
	go vet

lint:
	go get github.com/golang/lint/golint
	golint *.go

test:
	cd sysfs && rm -rf fixtures && tar xzf fixtures.tar.gz
	go test -v ./...

.PHONY: fmt lint test ci
