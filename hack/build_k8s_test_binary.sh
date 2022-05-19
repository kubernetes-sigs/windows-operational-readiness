#!/usr/bin/env bash

KUBERNETES_VERSION="v1.24.0"
KUBERNETES_HASH=6150737d11fa93a0b9ae4b32546a3ef96ab5dbe1

set -o errexit
set -o pipefail

if [ $1 == 1 ]; then
  if [ -d "./kubernetes" ]; then
    rm -fr "./kubernetes"
  fi
  git clone https://github.com/kubernetes/kubernetes.git
  pushd kubernetes
    git config --global --add safe.directory /home/kubo/op-readiness/kubernetes
    git checkout $KUBERNETES_HASH
    make WHAT="test/e2e/e2e.test"
  popd

  cp ./_output/bin/e2e.test ../e2e.test
  rm -rf kubernetes
else
  curl -L "https://dl.k8s.io/$KUBERNETES_VERSION/kubernetes-test-linux-amd64.tar.gz" -o /tmp/test.tar.gz
  tar xvzf /tmp/test.tar.gz --strip-components=3 kubernetes/test/bin/e2e.test
fi
