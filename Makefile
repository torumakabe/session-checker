APP=session-checker

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: build
build: clean
	go build -o ./bin/${APP} ./cmd/main.go

.PHONY: release-test
release-test:
	goreleaser --snapshot --skip-publish --rm-dist
