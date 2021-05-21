BINARY := $(shell basename "$(PWD)")
SOURCES := ./
GIT_COMMIT := $(shell git rev-list -1 HEAD)


.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.DEFAULT_GOAL := help

## build: Build the command line tool
build: clean
	CGO_ENABLED=0 go build \
	-ldflags '-w -extldflags "-static" -X main.gitCommit=$(GIT_COMMIT)' \
	-o ${BINARY} ${SOURCES}

## release: Build release artifacts
release:
	goreleaser --rm-dist --snapshot --skip-publish

## pack: Shrink the binary size
pack: build 
	upx -9 ${BINARY}

## test: Starts unit test
test:
	go test -v ./... -coverprofile coverage.out

## clean: Clean the binary
clean:
	rm -f $(BINARY)
