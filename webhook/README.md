## Setting up a Windows pod selector webhook

This document describes the steps to set up a mutating webhook to ensure successful scheduling of pods on a Windows
node. It achieves this by adding a `kubernetes.io/os: windows` key value node selector for each pod. The following steps
are derived based on the webhook maintained by the
Kubernetes [`sig-windows`](https://github.com/kubernetes/community/tree/master/sig-windows) team. For reference,
see [HyperV testing README](https://github.com/kubernetes-sigs/windows-testing/blob/master/helpers/hyper-v-mutating-webhook/README.md).

Steps

1. Create your cluster and ensure the cluster has at least one Windows node to test
2. Create a TLS certificate for secure communication
    - The Kubernetes API server uses HTTPS to communicate with webhooks. To support HTTPS, add a TLS certificate.
      You can either import an existing one from an externally trusted Certificate Authority (CA), or create your own
      using the following script and save the CA bundle for the next section.

```
./generate-cert.sh
```

3. Update the `caBundle` in `deployment.yaml` with the CA bundle you created the last step e.g.

```
...
   caBundle: |
     LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSURMekNDQWhlZ0F3SUJBZ0lVUnBOMHU0c3NqbFV5
     ...
     LS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
...
```

4. Update handler name in `runtimeclass.yaml` with the appropriate handler name for your container runtime e.g.

```
...
handler: runhcs-wcow-process
...
```

5. Run the `setup.sh` script to apply the YAML manifests. The script performs various actions including
    - taint the Windows nodes to ensure only pods with tolerations can schedule on a Windows node
    - store the TLS certificate and its corresponding private key as a secret

5. Run the `setup.sh` script to apply the YAML manifests. The script performs various actions including
    - taint the Windows nodes to ensure only pods with tolerations can schedule on a Windows node
    - store the TLS certificate and its corresponding private key as a secret
    - installing cert manager
    - installing the webhook

6. Verify that the webhook is correctly running
    - You should see the `windows-webhook-controller-manager deployment` in a running state

```
kubect get pods -n windows-webhook-system
```

7. Test if the webhook is working by deploying a pod without a node selector
    - Verify if the webhook took effect correctly by checking the logs. You should see a 200 code response with a log
      such as the following
   ```
   ...
   2023-10-24T06:29:00Z	INFO	webhook	Pod win-webserver-2019-56f7d48b87-ncmlz is being mutated
   2023-10-24T06:29:00Z	DEBUG	controller-runtime.webhook.webhooks	wrote response	{"webhook": "/mutate-v1-pod", "code": 200, "reason": "", "UID": "0c6a07f2-71a0-4d6d-a756-e75704da4d3e", "allowed": true}
   ...
   ```
    - Describe the pod to see if the node selector is automatically updated to the following
   ```
   ...
   Node-Selectors:   kubernetes.io/arch=amd64
                     kubernetes.io/os=windows
   ...
   ```
8. After verifying that the webhook is running and testing if the webhook is correctly taking effect, you are ready to
   begin running the Ops Readiness Windows focused test suite, with confidence that all pods will be correctly scheduled
   on a Windows node.