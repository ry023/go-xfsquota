.PHONY: test
test:
	go test ./... -coverprofile=coverage.out -covermode=count
