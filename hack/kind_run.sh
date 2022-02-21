set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# our exit handler (trap)
cleanup() {
    # always attempt to dump logs
    kind "export" logs "${ARTIFACTS}/logs" || true
    # KIND_IS_UP is true once we: kind create
    if [[ "${KIND_IS_UP:-}" = true ]]; then
        kind delete cluster || true
    fi
    unset IMG_REPO IMG_NAME IMG_TAG
}

# up a cluster with kind
create_cluster() {
    # create the config file
    cat <<EOF > "${ARTIFACTS}/kind-config.yaml"
# config for 1 control plane node and 2 workers
# necessary for conformance
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
# the control plane node
- role: control-plane
- role: worker
- role: worker
EOF

    # mark the cluster as up for cleanup
    # even if kind create fails, kind delete can clean up after it
    KIND_IS_UP=true

    KUBECONFIG="${HOME}/.kube/kind-config-default"
    export KUBECONFIG

    KIND_NODE_VERSION=v1.21.1
    export KIND_NODE_VERSION

    # actually create, with:
    # - do not delete created nodes from a failed cluster create (for debugging)
    # - wait up to one minute for the nodes to be "READY"
    # - use our multi node config
    kind create cluster \
        --image=kindest/node:${KIND_NODE_VERSION} \
        --retain \
        --wait=1m \
        "--config=${ARTIFACTS}/kind-config.yaml"
}


run_tests() {
    # create the op-readiness deployment
    cat <<EOF > "${ARTIFACTS}/op-readiness-deployment.yaml"
apiVersion: apps/v1
kind: Deployment
metadata:
  name: op-readiness
spec:
  replicas: 2
  selector:
    matchLabels:
      name: op-readiness
  template:
    metadata:
      labels:
        name: op-readiness
    spec:
      containers:
      - name: op-readiness
        image: ${IMG_REPO}/${IMG_NAME}:${IMG_TAG}
        imagePullPolicy: IfNotPresent
        command: [ "sleep" ]
        args: [ "infinity" ]
EOF
    kubectl apply -f ${ARTIFACTS}/op-readiness-deployment.yaml
    
    kubectl wait --for=condition=Ready --timeout=600s pod -l name=op-readiness
}

# setup kind, build kubernetes, create a cluster, run the e2es
main() {
    # ensure artifacts exists when not in CI
    ARTIFACTS="${ARTIFACTS:-${PWD}/_artifacts}"
    mkdir -p "${ARTIFACTS}"
    export ARTIFACTS
    export IMG_REPO=$1 IMG_NAME=$2 IMG_TAG=$3

    # now build an run the cluster and tests
    trap cleanup EXIT
    create_cluster
    run_tests
}

main $1 $2 $3
