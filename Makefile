ci: fmt lint test

fmt:
	! gofmt -l *.go | read nothing
	go vet

lint:
	go get github.com/golang/lint/golint
	golint *.go

test: sysfs/fixtures/.copied
	go test -v ./...

sysfs/fixtures/.copied:
	cd sysfs/fixtures.src; \
	find . \( -type f -o -type l \) -exec sh -c ' \
		np=../fixtures/$$(echo {} | sed "s/_@colon@_/:/g"); \
		nd=../fixtures/$$(dirname $$np); \
		mkdir -p $$nd; \
		cp -a {} $${np};' \
	\;
	touch $@

.PHONY: fmt lint test ci
