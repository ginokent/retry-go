SHELL := /bin/bash
GO_PROJECT := github.com/djeeno/retry-go
REPOSITORY_ROOT := ~/go/src/${GO_PROJECT}
VERSION := v0.0.2
REVISION := $(shell git rev-parse HEAD)
BUILD_DATE := $(shell TZ=UTC date +%Y%m%d%H%M%S)
GO_VERSION := $(shell go version)
BUILD_OPTS := -ldflags '-X "main.version=${VERSION}" -X "main.hash=${REVISION}" -X "main.builddate=${BUILD_DATE}" -X "main.goversion=${GO_VERSION}"'
OPEN_CMD := $(shell if command -v explorer.exe; then true; elif command -v open; then true; else echo echo; fi)

##
# targets
##
.PHONY: help
.DEFAULT_GOAL := help
help:  ## display help docs
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-40s\033[0m %s\n", $$1, $$2}'

test:  ## run test
	# run test
	go test -cover -race *.go

test-v:  ## run test
	# run test
	mkdir -p _test
	go test -cover -coverprofile=_test/cover.out -race -v *.go
	go tool cover -html=_test/cover.out -o _test/cover.html
	${OPEN_CMD} _test/cover.html

check-uncommitted:
	@if [ "`git diff; git diff --staged`" != "" ]; then\
		echo "Uncommitted changes. Execute the following command:";\
		echo "git commit -m 'release ${VERSION}'";\
		false;\
	fi

release: check-uncommitted test ## release as ${VERSION}
	# release ${VERSION}
	git tag -a "${VERSION}" -m "release ${VERSION}"
	git push
	git push --tags
