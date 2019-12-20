# Copyright 2018 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

.PHONY: usage

GO_REF ?= release-branch.go1.13
WEBSITE_HEAD := $(shell git rev-parse HEAD)
WEBSITE_CLEAN := $(shell (git status --porcelain | grep -q .) && echo dirty || echo clean)
ifeq ($(WEBSITE_CLEAN),clean)
	DOCKER_VERSION ?= $(WEBSITE_HEAD)
else
	DOCKER_VERSION ?= $(WEBSITE_HEAD)-dirty
endif
GCP_PROJECT := golang-org
GCP_SERVICE := default
DOCKER_TAG := gcr.io/$(GCP_PROJECT)/golangorg:$(DOCKER_VERSION)

usage:
	@echo "See Makefile and README.md"
	@exit 1

cloud-build: Dockerfile.prod
	gcloud builds submit \
		--project=$(GCP_PROJECT) \
		--config=cloudbuild.yaml \
		--substitutions=_GO_REF=$(GO_REF),_WEBSITE_HEAD=$(WEBSITE_HEAD),_WEBSITE_CLEAN=$(WEBSITE_CLEAN),_DOCKER_TAG=$(DOCKER_TAG) \
		../.. # source code

docker-build: Dockerfile.prod
	# NOTE(cbro): move up in directory to include entire tools repo.
	# NOTE(cbro): any changes made to this command must also be made in cloudbuild.yaml.
	cd ../..; docker build \
		-f=cmd/golangorg/Dockerfile.prod \
		--build-arg=GO_REF=$(GO_REF) \
		--build-arg=WEBSITE_HEAD=$(WEBSITE_HEAD) \
		--build-arg=WEBSITE_CLEAN=$(WEBSITE_CLEAN) \
		--build-arg=DOCKER_TAG=$(DOCKER_TAG) \
		--build-arg=BUILD_ENV=local \
		--tag=$(DOCKER_TAG) \
		.

docker-push: docker-build
	docker push $(DOCKER_TAG)

deploy:
	gcloud -q app deploy app.prod.yaml \
		--project=$(GCP_PROJECT) \
		--no-promote \
		--image-url=$(DOCKER_TAG)

get-latest-url:
	@gcloud app versions list \
		--service=$(GCP_SERVICE) \
		--project=$(GCP_PROJECT) \
		--sort-by=~version.createTime \
		--format='value(version.versionUrl)' \
		--limit 1 | cut -f1 # NOTE(cbro): gcloud prints out createTime as the second field.

get-latest-id:
	@gcloud app versions list \
		--service=$(GCP_SERVICE) \
		--project=$(GCP_PROJECT) \
		--sort-by=~version.createTime \
		--format='value(version.id)' \
		--limit 1 | cut -f1 # NOTE(cbro): gcloud prints out createTime as the second field.

regtest:
	go test -v \
		-regtest.host=$(shell make get-latest-url) \
		-run=Live

publish: regtest
	gcloud -q app services set-traffic $(GCP_SERVICE) \
		--splits=$(shell make get-latest-id)=1 \
		--project=$(GCP_PROJECT)

	@echo !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	@echo Stop and/or delete old versions:
	@echo "https://console.cloud.google.com/appengine/versions?project=$(GCP_PROJECT)&serviceId=$(GCP_SERVICE)&versionsquery=%255B%257B_22k_22_3A_22env_22_2C_22t_22_3A10_2C_22v_22_3A_22_5C_22FLEXIBLE_5C_22_22_2C_22s_22_3Atrue%257D%255D"
	@echo !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
