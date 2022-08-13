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

# Kubernetes build
IMG_REPO ?= <default_value_if_not_set_in_environment>
IMG_NAME ?= op-readiness
IMG_TAG ?= dev
KUBERNETES_HASH ?= 0

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
	docker run --rm -v $(PWD):/app -w /app -it golangci/golangci-lint golangci-lint run -v

## --------------------------------------
## Build
## --------------------------------------
##@ build:

.PHONY: docker-build
docker-build:  ## Build the Docker image
	docker build -t ${IMG_REPO}/${IMG_NAME}:${IMG_TAG} .
	docker push ${IMG_REPO}/${IMG_NAME}:${IMG_TAG}

.PHONY: build
build:  ## Build the binary using local golang
	./hack/build_k8s_test_binary.sh ${KUBERNETES_HASH}
	go build -o ./op-readiness .

### --------------------------------------
### Setup
### --------------------------------------
##@ setup:

.PHONY: local-kind-test
local-kind-test: docker-build  ## Run e2e tests with Kind, useful for development mode
	./hack/kind_run.sh ${IMG_REPO} ${IMG_NAME} ${IMG_TAG}

### --------------------------------------
### Testing
### --------------------------------------
##@ testing:

.PHONY: sonobuoy-plugin
sonobuoy-plugin:  ## Run the Sonobuoy plugin
	sonobuoy delete
	sonobuoy run --sonobuoy-image projects.registry.vmware.com/sonobuoy/sonobuoy:v0.56.3 --plugin sonobuoy-plugin.yaml --wait

sonobuoy-results:  ## Read Sonobuoy results
	rm -rf sonobuoy-results
	mkdir sonobuoy-results
	$(eval OUTPUT=$(shell sonobuoy retrieve))
	tar -xf $(OUTPUT) -C sonobuoy-results