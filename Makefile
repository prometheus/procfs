ci: fmt lint test

fmt:
	! gofmt -l *.go | read nothing
	go vet

lint:
	go get github.com/golang/lint/golint
	golint *.go

test:
	rm -rf sysfs/fixtures; \
	cd sysfs/fixtures.src; \
	find . \( -type f -o -type l \) -exec sh -c ' \
		np=../fixtures/$$(echo {} | sed "s/_@colon@_/:/g"); \
		nd=../fixtures/$$(dirname $$np); \
		mkdir -p $$nd; \
		cp -a {} $${np};' \
	\;
	go test -v ./...

.PHONY: fmt lint test ci
