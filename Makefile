OUT := xpb.exe
PKG := github.com/swoldemi/xpb
VERSION := $(shell git describe --always --long)

.PHONY: build
build:
	 GOOS=windows go build -v -i -o $(OUT) -a -installsuffix cgo -tags netgo -ldflags '-w -extldflags "-static" -X main.GitSHA=${VERSION}' *.go

.PHONY: test
test:
	go test -v -race -timeout 30s -count=1 -coverprofile=profile.out ./...

# Static code analysis tooling and checks
.PHONY: check
check:
	goimports -w -l -e .
	gofmt -s -w .
	golangci-lint run ./... \
		-E goconst \
		-E gocyclo \
		-E gosec  \
		-E gofmt \
		-E maligned \
		-E misspell \
		-E nakedret \
		-E unconvert \
		-E unparam \
		-E dupl
	goreportcard-cli -v -t 90

.PHONY: update
update:
	go get $(shell go list -f "{{if not (or .Main .Indirect)}}{{.Path}}{{end}}" -m all)
	go mod tidy
