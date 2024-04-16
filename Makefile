tidy:
	@go mod tidy
	@gofumpt -l -w .

test:
	@go test ./...

test-coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	@go vet -unsafeptr=false ./...

lint:
	@golangci-lint run

ci-test:
	@go test -race -vet=off ./...

format:
	@gofumpt -l .