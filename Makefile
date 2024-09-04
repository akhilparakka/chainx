build:
	@go build -o ./bin/chainx

run: build
	@./bin/chainx

test:
	@go test ./... -v