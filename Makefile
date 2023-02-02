SHELL = /bin/bash
APP := session-checker
GO_ENV := CGO_ENABLED=0
GO_BUILD_FLAGS := \
	-ldflags '-s -w' \
	-trimpath
IMAGE_NAME ?= session-checker
IMAGE_TAG ?= local

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test-short
test-short:
	go test -v -short ./...

.PHONY: test
test:
	docker compose build --no-cache
	docker compose up -d
	docker exec session-checker bash -c "cd /test && GOCACHE=/tmp go test -v ./..."

.PHONY: build
build: clean
	${GO_ENV} go build ${GO_BUILD_FLAGS} -o ./bin/${APP} ./cmd/*

.PHONY: release-test
release-test:
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: docker-build
docker-build:
	docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .

.PHONY: docker-image-scan
docker-image-scan: docker-build
	trivy image ${IMAGE_NAME}:${IMAGE_TAG}
