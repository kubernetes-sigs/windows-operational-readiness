#!/usr/bin/env bash

KUBERNETES_VERSION="v1.24.0"

set -o errexit
set -o pipefail

if [ $1 != 0 ]; then
  if [ -d "./kubernetes" ]; then
    rm -fr "./kubernetes"
  fi
  git clone https://github.com/kubernetes/kubernetes.git
  pushd kubernetes
    git config --global --add safe.directory /home/kubo/op-readiness/kubernetes
    git checkout $1
    make WHAT="test/e2e/e2e.test"
  popd

  cp ./kubernetes/_output/bin/e2e.test ./e2e.test
  rm -rf kubernetes
else
  curl -L "https://dl.k8s.io/$KUBERNETES_VERSION/kubernetes-test-linux-amd64.tar.gz" -o /tmp/test.tar.gz
  tar xvzf /tmp/test.tar.gz --strip-components=3 kubernetes/test/bin/e2e.test
fi
