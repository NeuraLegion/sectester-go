build:
	find . -name go.mod -execdir go build ./... \;

clean:
	go clean

test:
	find . -name go.mod -execdir go test ./... \;

coverage:
	find . -name go.mod -execdir go test ./... -coverprofile=coverage.out \;

deps:
	find . -name go.mod -execdir go mod download \;

vet:
	find . -name go.mod -execdir go vet ./... \;

format:
	find . -name go.mod -execdir go fmt ./... \;

tidy:
	find . -name go.mod -execdir go mod tidy \;

lint:
	find . -name go.mod -execdir golangci-lint run \;