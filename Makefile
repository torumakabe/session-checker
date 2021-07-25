APP=session-checker

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: test
test:
	go test -v -short ./...

.PHONY: build
build: clean
	go build -o ./bin/${APP} ./cmd/*

.PHONY: release-test
release-test:
	goreleaser --snapshot --skip-publish --rm-dist
