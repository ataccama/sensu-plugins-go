.PHONY: all checks
all: checks

CHECK_SOURCES = $(patsubst %/main.go,%,$(shell find check* -mindepth 1 -maxdepth 1 -name '*.go'))
CHECK_TARGETS = $(patsubst %,../bin/%,$(CHECK_SOURCES))

checks: $(CHECK_TARGETS)

../bin/%: %/*.go
	go fmt $^
	go build -o $@ $^

