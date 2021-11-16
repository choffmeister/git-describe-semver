.PHONY: *

MAIN := .
TEST := . ./internal

run:
	go run $(MAIN)

test:
	go test -v $(TEST)

test-watch:
	watch -n1 go test -v $(TEST)

test-cover:
	go test -coverprofile=coverage.out $(TEST)
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

build:
	goreleaser build --rm-dist --snapshot

release:
	goreleaser release --rm-dist
