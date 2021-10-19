lint-install:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin

lint:
	golangci-lint run

test:
	go clean -testcache
	go test ./... --race

test-coverage:
	go clean -testcache
	go test ./... --race -coverprofile cover.out
	go tool cover -func cover.out
	@rm cover.out