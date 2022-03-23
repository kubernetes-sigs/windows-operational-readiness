IMG_REPO ?= <default_value_if_not_set_in_environment>
IMG_NAME ?= op-readiness
IMG_TAG ?= dev

.PHONY: local-kind-test
local-kind-test: docker-build
	./hack/kind_run.sh ${IMG_REPO} ${IMG_NAME} ${IMG_TAG}

.PHONY: docker-build
docker-build:
	docker build -t ${IMG_REPO}/${IMG_NAME}:${IMG_TAG} .
	docker push ${IMG_REPO}/${IMG_NAME}:${IMG_TAG}

.PHONY: build
build: 
	./hack/build_k8s_test_binary.sh
	go build -o ./op-readiness .

.PHONY: sonobuoy-plugin
sonobuoy-plugin:
	sonobuoy delete
	sonobuoy run --plugin sonobuoy-plugin.yaml

