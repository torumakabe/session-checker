APP=session-checker

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: test-short
test-short:
	go test -v -short ./...

.PHONY: test
test:
	docker-compose up -d
	docker exec session-checker bash -c "cd /test && go test -v ./..."

.PHONY: build
build: clean
	go build -o ./bin/${APP} ./cmd/*

.PHONY: release-test
release-test:
	goreleaser --snapshot --skip-publish --rm-dist
