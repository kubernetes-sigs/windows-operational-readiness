#!/bin/bash

windows_nodes=$(kubectl get nodes -l kubernetes.io/os=windows -o jsonpath='{.items[*].metadata.name}')
for node in $windows_nodes; do
  echo "Tainting Windows node $node"
  kubectl taint node "$node" os=windows:NoSchedule --overwrite
done

echo "Applying runtimeclass.yaml"
kubectl apply -f runtimeclass.yaml

echo "Installing cert manager"
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml

kubectl apply -f namespace.yaml

echo
echo ">> Generating kube secrets..."
kubectl create secret tls webhook-server-cert \
  -n windows-webhook-system \
  --cert=server.crt \
  --key=server.key

echo "Installing webhook deployment"
kubectl apply -f deployment.yaml

echo "Untainting Windows nodes"
for node in $windows_nodes; do
  kubectl taint node "$node" os=windows:NoSchedule- --overwrite
done