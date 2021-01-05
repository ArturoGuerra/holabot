.PHONY: all clean build docker docker-build docker-push

GORUN = go run
GOBUILD = go build

all: build

clean:
	rm -rf bin

build: clean
	$(GOBUILD) -o bin/hola cmd/hola/main.go

docker-build:
	test $(DOCKERREPO)
	docker build . -t $(DOCKERREPO)

docker-push:
	test $(DOCKERREPO)
	docker push $(DOCKERREPO)

docker: docker-build docker-push