.DEFAULT_GOAL := run

fmt:
	go fmt ./...
.PHONY:fmt

# go install golang.org/x/lint/golint@latest
lint: fmt
	@GOPATH=$(shell go env | grep GOPATH | awk -F= '{print $$2}' | tr -d '"'); \
	$$GOPATH/bin/golint ./...
.PHONY:lint

vet: fmt
	go vet ./...
.PHONY:vet

build: vet
	go build -o ./application_metadata_api_server
.PHONY:build

# run a server at :8080
run: vet
	go run main.go
.PHONY:run

# run unit tests
GOTARGET = .
TEST_PKGS ?= $(GOTARGET)/cache/... $(GOTARGET)/server/...

test:
	go test -v $(TEST_PKGS)
.PHONY:test

# run curl with sample data to insert and query
test-query:
	curl --data-binary  "@testdata/valid-payload1.yaml" http://localhost:8080/put
	@echo
	curl --data-binary  "@testdata/valid-payload2.yaml" http://localhost:8080/put
	@echo
	curl --data-binary  "1" http://localhost:8080/get
	@echo
	curl --data-binary  "@testdata/invalid-payload1.yaml" http://localhost:8080/put
	@echo
	curl --data-binary  "@testdata/invalid-payload2.yaml" http://localhost:8080/put
	@echo
	curl --data-binary  "@testdata/query1.yaml" 		http://localhost:8080/query
	@echo
	curl --data-binary  "@testdata/query2.yaml"			http://localhost:8080/query
	@echo
	curl --data-binary  "@testdata/query3.yaml" 		http://localhost:8080/query
	@echo
	curl --data-binary  "@testdata/query4.yaml" 		http://localhost:8080/query
	@echo
.PHONY:test-query

# build a docker image with the api server
docker-build:
	docker build --tag application_metadata_api_server .
.PHONY:docker-build

# run a server at :8080 inside a docker container
docker-run-server:
	docker run --rm -p 8080:8080 application_metadata_api_server
.PHONY:docker-run-server

docker-run-go-compiler-shell:
#-t Allocate a pseudo-TTY
#-v, --volume list                    Bind mount a volume
#-i, --interactive                    Keep STDIN open even if not attached
#--rm                             Automatically remove the container when it exits
	docker run -t --rm -i -v `pwd`:/application_metadata_api_server  golang:1.18-buster /bin/bash
.PHONY:docker-run-go-compiler-shell
