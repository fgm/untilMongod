all: lint install

.phony: lint
lint:
	go tool staticcheck ./...


.phony: install
install:
	go install .