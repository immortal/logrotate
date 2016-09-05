.PHONY: all get test cover

GO ?= go

all: test

get:
	${GO} get

test: get
	${GO} test -race -v

cover:
	${GO} test -cover && \
	${GO} test -coverprofile=coverage.out  && \
	${GO} tool cover -html=coverage.out
