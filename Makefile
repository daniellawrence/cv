
BUF=buf

.PHONY: setup proto test clean


setup:
	go install github.com/bufbuild/buf/cmd/buf@latest

proto:
	cd shared/proto && $(BUF) generate

test:
	cd backend/search && go test -v ./...