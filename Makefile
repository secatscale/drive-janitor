.PHONY: test format

test:
	go test ./... -v

format: 
	gofmt -w .
