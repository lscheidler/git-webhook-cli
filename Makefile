all: fmt vet build

build:
	go build

fmt:
	go fmt $(shell find . -type d | egrep -v "(vendor|.git)")

vet:
	go vet
