# Copyright 2022 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHELL := /usr/bin/env bash
.DEFAULT_GOAL := help

# Container build
STAGING_REGISTRY ?= gcr.io/k8s-staging-win-op-rdnss
IMG_NAME ?= k8s-win-op-rdnss
TAG ?= $(shell git describe --tags --always `git rev-parse HEAD`)
IMG_PATH ?= $(STAGING_REGISTRY)/$(IMG_NAME)

# Kubernetes version
KUBERNETES_HASH ?= 0
KUBERNETES_VERSION ?= v1.24.0

## --------------------------------------
## Help
## --------------------------------------
##@ help:

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

## --------------------------------------
## Linters
## --------------------------------------
##@ lint:

.PHONY: lint-go
lint-go: ## Lint codebase
	docker run --rm -v $(PWD):/app -w /app -it golangci/golangci-lint golangci-lint run -v --fix

## --------------------------------------
## Clean
## --------------------------------------
##@ clean:

.PHONY: clean
clean :
	rm -rf op-readiness test.tar.gz e2e.test ./--report-prefix/

## --------------------------------------
## Build
## --------------------------------------
##@ build:

.PHONY: build
build:  ## Build the binary using local golang
	./hack/build_k8s_test_binary.sh ${KUBERNETES_HASH}
	go build -o ./op-readiness .

## --------------------------------------
## Container
## --------------------------------------
##@ container:

.PHONY: image-build
image-build: ## Build the container image
	docker build --build-arg KUBERNETES_VERSION=$(KUBERNETES_VERSION) -t $(IMG_PATH):$(TAG) .
	docker tag $(IMG_PATH):$(TAG) $(IMG_PATH):latest

.PHONY: image_push
image-push: ## Push the container image to k8s-staging bucket
	docker push $(IMG_PATH):$(TAG)
	docker push $(IMG_PATH):latest

.PHONY: release-staging
release-staging: ## Builds and push container image to k8s-staging bucket
	$(MAKE) image-build image-push

### --------------------------------------
### Setup
### --------------------------------------
##@ setup:

.PHONY: local-kind-test
local-kind-test: image-build ## Run e2e tests with Kind, useful for development mode
	./hack/kind_run.sh ${IMG_REPO} ${IMG_NAME} ${TAG}

### --------------------------------------
### Testing
### --------------------------------------
##@ testing:

.PHONY: sonobuoy-plugin
sonobuoy-plugin:  ## Run the Sonobuoy plugin
	sonobuoy delete --all
	sonobuoy run --sonobuoy-image projects.registry.vmware.com/sonobuoy/sonobuoy:v0.56.9 --plugin sonobuoy-plugin.yaml --wait

.PHONY: sonobuoy-results
sonobuoy-results:  ## Read Sonobuoy results
	$(eval OUTPUT=$(shell sonobuoy retrieve))
	sonobuoy results --mode=report $(OUTPUT)

.PHONY: sonobuoy-config-gen
sonobuoy-config-gen:  ## Run the Sonobuoy plugin
	cd sonobuoy; ytt -f . > ../sonobuoy-plugin.yaml
