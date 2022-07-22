IMG_REPO ?= <default_value_if_not_set_in_environment>
IMG_NAME ?= op-readiness
IMG_TAG ?= dev
KUBERNETES_HASH ?= 0

.PHONY: local-kind-test
local-kind-test: docker-build
	./hack/kind_run.sh ${IMG_REPO} ${IMG_NAME} ${IMG_TAG}

.PHONY: docker-build
docker-build:
	docker build -t ${IMG_REPO}/${IMG_NAME}:${IMG_TAG} .
	docker push ${IMG_REPO}/${IMG_NAME}:${IMG_TAG}

.PHONY: build
build: 
	./hack/build_k8s_test_binary.sh ${KUBERNETES_HASH}
	go build -o ./op-readiness .

.PHONY: sonobuoy-plugin
sonobuoy-plugin:
	sonobuoy delete
	sonobuoy run --sonobuoy-image projects.registry.vmware.com/sonobuoy/sonobuoy:v0.56.3 --plugin sonobuoy-plugin.yaml --wait

sonobuoy-results:
	rm -rf sonobuoy-results
	mkdir sonobuoy-results
	$(eval OUTPUT=$(shell sonobuoy retrieve))
	tar -xf $(OUTPUT) -C sonobuoy-results