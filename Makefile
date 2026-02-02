.PHONY: build
build:
	go build -o uss main.go

.PHONY: test
test:
	go fmt ./...
	go mod tidy
	go test ./...

.PHONY: gen-golden
gen-golden:
	cp gen_builtin/output/output.sh builtin_complete.sh
	cp gen_builtin/output/output_zsh.sh zsh_builtin_complete.sh
