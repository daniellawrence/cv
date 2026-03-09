
BUF=buf

.PHONY: setup proto test clean


setup:
	go install github.com/bufbuild/buf/cmd/buf@latest

proto: setup
	cd proto && $(BUF) generate

test: proto
	cd backend/education && go test ./...