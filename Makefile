BINARY_NAME=reservr

.PHONY: build
build:
	GOARCH=amd64 GOOS=darwin go build -o build/${BINARY_NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o build/${BINARY_NAME}-linux main.go

run:
	go run main.go

clean:
	go clean
	rm build/${BINARY_NAME}-darwin
	rm build/${BINARY_NAME}-linux

test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet

lint:
	golangci-lint run --enable-all