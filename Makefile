
BUF=buf

.PHONY: setup proto test clean


setup:
	go install github.com/bufbuild/buf/cmd/buf@latest

proto: setup
	cd proto && $(BUF) generate

test: proto
	cd backend/education && go test ./...

bin/golangci-lint:
	curl -sSfL https://golangci-lint.run/install.sh | sh -s v2.11.2

lint:
	find backend -mindepth 1  -type d -exec ./bin/golangci-lint run {} \;