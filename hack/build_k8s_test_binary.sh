#!/usr/bin/env bash

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

set -o errexit
set -o pipefail
set -x

# todo(knabben) - fetch latest or pass as argument
KUBERNETES_VERSION=${KUBERNETES_VERSION:-"v1.24.0"}
KUBERNETES_REPO=${KUBERNETES_REPO:-"https://github.com/kubernetes/kubernetes.git"}

if [ $1 != 0 ]; then
  # Using the hash passed as argument
  if [ -d "./kubernetes" ]; then
    rm -fr "./kubernetes"
  fi

  git clone ${KUBERNETES_REPO}
  pushd kubernetes
    git config --global --add safe.directory ${PWD}
    git checkout $1
    make WHAT="test/e2e/e2e.test"
  popd

  # Copy the binary
  cp ./kubernetes/_output/bin/e2e.test ./e2e.test
  # Clean up folder
  rm -rf kubernetes
elif [ ! -f "e2e.test" ]; then
  # Download the binary directly from Kubernetes release, skip if already exists
  curl -L "https://dl.k8s.io/${KUBERNETES_VERSION}/kubernetes-test-linux-amd64.tar.gz" -o test.tar.gz
  tar xvzf test.tar.gz --strip-components=3 kubernetes/test/bin/e2e.test
fi
