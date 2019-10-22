all: fmt build

build:
	go build

fmt:
	go fmt $(shell find . -type d | egrep -v "(vendor|.git)")
