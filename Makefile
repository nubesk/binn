PORT_SRC = 8080
PORT_DST = 8080
MOUNT_PATH_SRC = $(PWD)
MOUNT_PATH_DST = /go/src/github.com/binn

build-image:
	docker build -t binn -f Dockerfile .

build-image-dev:
	docker build -t binn-dev -f Dockerfile.dev .

run:
	docker run -it --rm -p $(PORT_SRC):$(PORT_DST) binn

run-dev:
	docker run -it --rm -p $(PORT_SRC):$(PORT_DST) -v $(MOUNT_PATH_SRC):$(MOUNT_PATH_DST) binn-dev /bin/sh

test:
	go test ./binn ./server
