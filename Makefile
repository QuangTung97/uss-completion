.PHONY: build
build:
	go build -o uss main.go

.PHONY: test
test:
	go test ./...
