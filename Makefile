.PHONY: all get test clean cover

GO ?= go

all: test

clean:
	@rm -rf *.out

get:
	${GO} get

test: get
	${GO} test -race -v

cover:
	${GO} test -cover && \
	${GO} test -coverprofile=coverage.out  && \
	${GO} tool cover -html=coverage.out
