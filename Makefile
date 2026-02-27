.PHONY: build
build:
	go build -o bin/fixvars ./cmd/fixvars

.PHONY: test
test:
	go test
