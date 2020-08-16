.PHONY: run test test-watch build push

# The binary to build (just the basename).
BIN := git-describe-semver

# Where to push the docker image.
REGISTRY ?= choffmeister

IMAGE := $(REGISTRY)/$(BIN)

# This version-strategy uses git tags to set the version string
VERSION := $(shell git describe --tags --always --dirty)

run:
	go run .

run-watch:
	watch -n1 go run .

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o git-describe-semver-linux-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o git-describe-semver-darwin-amd64 .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o git-describe-semver-windows-amd64 .

test:
	go test -v

test-watch:
	watch -n1 go test -v

test-cover:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

container:
	docker build -t $(IMAGE):$(VERSION) .

container-push: container
	docker push $(IMAGE):$(VERSION)
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
	docker push $(IMAGE):latest
