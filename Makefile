
BUF=buf

.PHONY: setup proto test clean


setup:
	go install github.com/bufbuild/buf/cmd/buf@latest

proto:
	cd proto && $(BUF) generate

test:
	cd backend/education && go test ./...