SHELL = /bin/bash

TESTS_ALLOWED=TRUE
HELLOWORLD_TESTS=TRUE

CMD_DIRS := \
	./cmd/helloworld

SRC_DIRS := $(shell \
	find . -name "*.go" -not -path "./vendor/*" | \
	xargs -I {} dirname {}  | \
	uniq)

TEST_DIRS := $(shell \
	find . -name "*_test.go" -not -path "./vendor/*" | \
	xargs -I {} dirname {}  | \
	uniq)

TOOLS := \
	github.com/golang/dep/cmd/dep; \
	golang.org/x/lint/golint; \
	golang.org/x/tools/cmd/goimports; \
	github.com/stretchr/testify;

GETTOOLS := $(foreach TOOL,$(TOOLS),go get -u $(TOOL))

default: tools deps fmt vet lint build

all: check tools deps fmt vet lint test build

prod:
	go install -v $(CMD_DIRS)

install: check tools deps fmt vet lint test
	go install -v $(CMD_DIRS)

check:
	@[ "${TESTS_ALLOWED}" ] || ( echo ">> HELLOWORLD_TESTS is not set"; exit 1 )
	@[ "${TESTS_ALLOWED}" == TRUE ] || ( echo ">> HELLOWORLD_TESTS is not set to TRUE"; exit 1 )

tools:
	$(GETTOOLS)

deps:
	dep ensure

fmt: $(SRC_DIRS)
	@for dir in $^; do \
		pushd ./$$dir > /dev/null ; \
		goimports -w -local "github.com/go-daemons/helloworld" *.go ; \
		popd > /dev/null ; \
	done;

vet: $(SRC_DIRS)
	@for dir in $^; do \
		pushd ./$$dir > /dev/null ; \
		go vet ; \
		popd > /dev/null ; \
	done;

lint: $(SRC_DIRS)
	@for dir in $^; do \
		pushd ./$$dir > /dev/null ; \
		golint -set_exit_status ; \
		popd > /dev/null ; \
	done;

test: $(TEST_DIRS)
	@[ "${TESTS_ALLOWED}" ] || ( echo ">> HELLOWORLD_TESTS is not set"; exit 1 )
	@[ "${TESTS_ALLOWED}" == TRUE ] || ( echo ">> HELLOWORLD_TESTS is not set to TRUE"; exit 1 )
	@for dir in $^; do \
		pushd ./$$dir > /dev/null ; \
		go test -v ; \
		popd > /dev/null ; \
	done;

build: $(CMD_DIRS)
	@for dir in $^; do \
		pushd ./$$dir > /dev/null ; \
		go build -v ; \
		popd > /dev/null ; \
	done;

clean:
	go clean $(CMD_DIRS)

uninstall:
	go clean -i $(CMD_DIRS)

.PHONY: check prod all tools dep fmt vet lint test install build clean
