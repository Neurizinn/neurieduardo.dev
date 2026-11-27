.PHONY: setup generate fmt build dev test check vercel-output
setup:
	go mod download
generate:
	go run github.com/a-h/templ/cmd/templ@v0.3.960 generate
fmt:
	go run github.com/a-h/templ/cmd/templ@v0.3.960 fmt .
	go fmt ./...
build: generate
	go run ./cmd/build
dev: build
	go run ./cmd/serve
test:
	go test ./...
check: generate test build
	go vet ./...
vercel-output: build
	go run ./cmd/vercel-output
