all: test

TEST_ARGS ?= -v -count=1 -timeout 120s

test:
	go test $(TEST_ARGS) ./...

test.debug:
	go test -tags debug $(TEST_ARGS) ./...
