build:
	go build ./...

clean:
	go clean

test:
	go test ./...

coverage:
	go test ./... -coverprofile=coverage.out

deps:
	go mod download

vet:
	go vet ./...

format:
	go fmt ./...

tidy:
	go mod tidy

lint:
	golangci-lint run
