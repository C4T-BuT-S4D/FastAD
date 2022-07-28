.PHONY: lint-proto
lint-proto:
	cd proto && buf lint && buf build

.PHONY: lint-go
lint-go:
	golangci-lint run -v -c .golangci.yml ./...

.PHONY: lint-front
lint-front:
	cd front && yarn lint

.PHONY: proto
proto: lint-proto
	cd proto && buf generate

.PHONY: tidy
tidy:
	go mod tidy -compat="1.18"

.PHONY: lint
lint: lint-proto lint-front lint-go

.PHONY: test-go
test-go:
	go test -race ./...

.PHONY: test
test: test-go
