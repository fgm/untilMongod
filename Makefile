all: lint install

GO       := go
PACKAGES := ./..

.PHONY: lint cover test build install

lint:
	$(GO) tool staticcheck $(PACKAGES)

cover:
	# This runs the benchmarks just once, as unit tests, for coverage reporting only.
	# It does not replace running "make bench".
	$(GO) test -v -race -run=. -coverprofile=coverage/cover.out -covermode=atomic ./...

test:
	# This includes the fuzz tests in unit test mode
	$(GO) test -race $(PACKAGES)

build: test
	$(GO) build $(PACKAGES)

install: build
	$(GO) install .