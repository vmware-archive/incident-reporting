# Copyright 2019 VMware, Inc.
# SPDX-License-Identifier: BSD-2

# tag is the git tag or diff from the tag
TAG := $(shell git --no-pager describe --tags --always)

# or if empty, tag is the git revision
ifeq ($(TAG),)
TAG := $(shell git rev-parse --verify HEAD)
endif

# or if empty, tag is latest
ifeq ($(TAG),)
TAG := latest
endif

BASEREPO   ?= index.docker.io/tompscanlan
UIIMG      := incident-reporting-ui:${TAG}
TRUFFLEIMG := incident-reporting-truffle:${TAG}
LATESTTAG  := latest

all: containers
containers: ui truffle
push: push-ui push-truffle
.PHONY: ui truffle push-ui push-truffle containers

ui: ui-server/*.go ui-server/Dockerfile truffle/contracts/*.sol
	docker build -t $(UIIMG) -f ui-server/Dockerfile .
	docker tag $(UIIMG) $(BASEREPO)/$(UIIMG)

truffle: truffle/package*.json truffle/Dockerfile truffle/contracts/*.sol truffle/migrations/*.js truffle/test/*.js truffle/*.js
	docker build -t $(TRUFFLEIMG) -f truffle/Dockerfile .
	docker tag $(TRUFFLEIMG) $(BASEREPO)/$(TRUFFLEIMG)

push-ui: ui
	docker push $(BASEREPO)/$(UIIMG)

push-truffle: truffle
	docker push $(BASEREPO)/$(TRUFFLEIMG)

run-truffle:
	docker run -e PRODUCTION_URL -it incident-reporting-truffle:$(TAG) bash -s 'echo -e "#run this:\ntruffle deploy --network production --reset" && bash'

run-ui:
	docker run -it \
		-v ${PWD}/keystore:/keystore \
		--env-file env \
		-p 8080:80 \
		incident-reporting-ui:$(TAG)
