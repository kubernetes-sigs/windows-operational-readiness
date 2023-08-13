DOCKER_IMAGE="docker.io/library/golang:1.21.0"
REPO="docker.io/winopreadiness/rextension:dev"

.PHONY=build
build:
	docker build --build-arg builder_image=${DOCKER_IMAGE} \
	--build-arg ARCH=amd64 -t ${REPO} \
	--file ./Dockerfile .
	docker push ${REPO}
