IMG_REPO ?= <default_value_if_not_set_in_environment>
IMG_NAME ?= op-readiness
IMG_TAG ?= dev

.PHONY: local-kind-test
local-kind-test: docker-build
	./kind_run.sh ${IMG_REPO} ${IMG_NAME} ${IMG_TAG}

.PHONY: docker-build
docker-build:
	docker build -t ${IMG_REPO}/${IMG_NAME}:${IMG_TAG} .
	docker push ${IMG_REPO}/${IMG_NAME}:${IMG_TAG}