.PHONY: test format all re clean

all:
	@echo "Building the project..."
	go build ./main.go
	@echo "Build complete"

re: clean all

clean:
	@echo "Cleaning up..."
	rm -f drive-janitor
	@echo "Clean complete"

test:
	go test ./... -v

format: 
	gofmt -w .
