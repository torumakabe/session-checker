APP := session-checker
GO_ENV := CGO_ENABLED=0
GO_BUILD_FLAGS := \
	-ldflags '-s -w' \
	-trimpath

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: test-short
test-short:
	go test -v -short ./...

.PHONY: test
test:
	docker compose up -d
	docker exec session-checker bash -c "cd /test && go test -v ./..."

.PHONY: build
build: clean
	${GO_ENV} go build ${GO_BUILD_FLAGS} -o ./bin/${APP} ./cmd/*

.PHONY: release-test
release-test:
	goreleaser --snapshot --skip-publish --rm-dist
