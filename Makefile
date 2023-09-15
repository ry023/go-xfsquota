.PHONY: test
test:
	go test ./... -coverprofile=coverage.out -covermode=count

depsdev:
	go install github.com/Songmu/gocredits/cmd/gocredits@latest

prerelease_for_tagpr: depsdev
	go mod tidy
	gocredits -w .
	git add CHANGELOG.md CREDITS go.mod go.sum
