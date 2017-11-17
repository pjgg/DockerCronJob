
SRC_PATH=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION=v0.0.1

clean:
	rm -rf\
	 $(SRC_PATH)/dist\
	 $(SRC_PATH)/debug\

all: clean sync version compile build_docker push_docker

compile:
	 $(SRC_PATH)/build/build-go.sh

build_docker:
	$(SRC_PATH)/build/build-docker.sh

push_docker:
	$(SRC_PATH)/build/docker-push.sh

version:
	@echo $(VERSION)	

sync:
	- go get -v github.com/kardianos/govendor
	- cd $(SRC_PATH) && govendor sync

