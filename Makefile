.PHONY: test
test: build_fake
	go test ./... -coverprofile=coverage.out -covermode=count

lint:
	golangci-lint run ./...

depsdev:
	go install github.com/Songmu/gocredits/cmd/gocredits@latest

build_fake:
	go build -o fake_xfs_quota testutil/cmd/fake_xfs_quota/main.go

prerelease_for_tagpr: depsdev
	go mod tidy
	gocredits -w .
	git add CHANGELOG.md CREDITS go.mod go.sum
