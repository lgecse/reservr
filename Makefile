BINARY_NAME=reservr
REPO=481155165509.dkr.ecr.eu-west-1.amazonaws.com
IMG=reservr-lambda

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

debug-run:
	docker run --platform linux/amd64 -d -v ~/.aws-lambda-rie:/aws-lambda -p 8080:8080 --entrypoint /aws-lambda/aws-lambda-rie kristofgyuracz/reservr:0.0.1 ./reservr-server

docker-build:
	docker build --platform linux/amd64 -t ${REPO}/${IMG} .

docker-push:
	docker push ${REPO}/${IMG}

docker-login:
	aws ecr get-login-password --region=eu-west-1 | docker login --username AWS --password-stdin ${REPO}