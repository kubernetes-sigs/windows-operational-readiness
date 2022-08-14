# Operational Readiness Runtime Hook

The propose of this experiment is to run the Sonobuoy test suite after the Workload control
plane get ready, meaning we can automate the run of the plugin on new clusters by default using
ClusterAPI and ClusterClass. The plugin only uses the lifecycle hook of AfterControlPlaneInitialized and
retry until the CP gets ready, the official [docs](https://cluster-api.sigs.k8s.io/tasks/experimental-features/runtime-sdk/implement-lifecycle-hooks.html)
can give more details about the lifecycle hooks and proper documentation.

*NOTE: This is a very experimental way to run this Sonobouy plugin. Don't try it on production.*
*NOTE: Make sure the Workload cluster has a CNI installed, on this setup, the CNI must be installed manually.*

### Development mode

The container and webhook is not published, and [Tilt](tilt.dev) is being used for the development
and run, this plugin was tested only on CAPD, but should work on other providers.

For the directories in the project: `config` is responsible to render all the objects on your clusters
`certmanager` issues a certificate for your webhook, `default` has the extensions deployments and RBAC
parsed by `kustomize`.

The `Dockerfile` and `tilt-addon.yaml` provides the resource for 
`extensionconfig.yaml` is the ExtensionConfig object, and installs the hooks and shows the
CAPI controller what service to access it. The other golang files are the source
code of the project, given the logic resides inside `handlers/handlers.go`.

```shell
├── config
│   ├── certmanager
│   │   ├── certificate.yaml
│   │   ├── kustomization.yaml
│   │   └── kustomizeconfig.yaml
│   └── default
│       ├── extension_image_patch.yaml
│       ├── extension_pull_policy.yaml
│       ├── extension_webhook_patch.yaml
│       ├── extension.yaml
│       ├── kustomization.yaml
│       ├── kustomizeconfig.yaml
│       ├── rolebinding.yaml
│       ├── role.yaml
│       ├── service_account.yaml
│       └── service.yaml
├── Dockerfile
├── extensionconfig.yaml
├── go.mod
├── go.sum
├── handlers
│   └── handlers.go
├── main.go
├── Makefile
├── README.md
└── tilt-addon.yaml
```

### Running the workload

It's possible to run `tilt up` directly from the folder, make sure the management cluster was initialized
with `${CAPI_HOME}/hack/kind-install-for-capd.sh` for example. Observe a new pod for the webserver in
the default namespace.

```shell
NAME                             READY   STATUS    RESTARTS   AGE
win-extension-5fbcd76b58-jvftz   1/1     Running   0          31m
```

Looking the logs it's possible to visualize the initialization of the plugin in the remote cluster,
the one can notice the retry of the extension until the CP gets ready, and the sonobuoy installation 
in sequence.

```shell
I0820 19:47:59.247722      63 logr.go:249] "setup: Starting manager"
I0820 19:47:59.247834      63 logr.go:249] "setup: Starting Runtime Extension server"
I0820 19:47:59.247942      63 internal.go:362] "Starting server" path="/metrics" kind="metrics" addr="[::]:8080"
I0820 19:47:59.248058      63 controller.go:185] "Starting EventSource" controller="remote/clustercache" controllerGroup="cluster.x-k8s.io" controllerKind="Cluster" source="kind source: *v1beta1.Cluster"
I0820 19:47:59.248086      63 controller.go:193] "Starting Controller" controller="remote/clustercache" controllerGroup="cluster.x-k8s.io" controllerKind="Cluster"
I0820 19:47:59.248244      63 server.go:148] "controller-runtime/webhook: Registering webhook" path="/hooks.runtime.cluster.x-k8s.io/v1alpha1/aftercontrolplaneinitialized/after-controlplane-initialized"
I0820 19:47:59.248390      63 server.go:148] "controller-runtime/webhook: Registering webhook" path="/hooks.runtime.cluster.x-k8s.io/v1alpha1/discovery"
I0820 19:47:59.248452      63 server.go:216] "controller-runtime/webhook/webhooks: Starting webhook server"
I0820 19:47:59.248744      63 logr.go:249] "controller-runtime/certwatcher: Updated current TLS certificate"
I0820 19:47:59.248867      63 logr.go:249] "controller-runtime/webhook: Serving webhook server" host="" port=9443
I0820 19:47:59.248977      63 logr.go:249] "controller-runtime/certwatcher: Starting certificate watcher"
I0820 19:47:59.349381      63 controller.go:227] "Starting workers" controller="remote/clustercache" controllerGroup="cluster.x-k8s.io" controllerKind="Cluster" worker count=10
I0820 19:48:38.749097      63 handlers.go:51] "AfterControlPlaneInitialized is called, trying to start the Sonobuoy plugin in the WL cluster."
E0820 19:48:38.749161      63 handlers.go:55] "Cluster Workload Control Plane is not ready yet, retrying."
...
I0820 19:48:43.046687      63 cluster_cache_tracker.go:238] "Creating cluster accessor for cluster \"default/development-2895\" with the regular apiserver endpoint \"https://172.18.0.6:6443\""
time="2022-08-20T19:48:43Z" level=info msg="create request issued" name=sonobuoy namespace= resource=namespaces
time="2022-08-20T19:48:43Z" level=info msg="create request issued" name=sonobuoy-serviceaccount namespace=sonobuoy resource=serviceaccounts
time="2022-08-20T19:48:43Z" level=info msg="create request issued" name=sonobuoy-serviceaccount-sonobuoy namespace= resource=clusterrolebindings
time="2022-08-20T19:48:43Z" level=info msg="create request issued" name=sonobuoy-serviceaccount-sonobuoy namespace= resource=clusterroles
time="2022-08-20T19:48:43Z" level=info msg="create request issued" name=sonobuoy-config-cm namespace=sonobuoy resource=configmaps
time="2022-08-20T19:48:43Z" level=info msg="create request issued" name=sonobuoy-plugins-cm namespace=sonobuoy resource=configmaps
time="2022-08-20T19:48:43Z" level=info msg="create request issued" name=sonobuoy namespace=sonobuoy resource=pods
time="2022-08-20T19:48:43Z" level=info msg="create request issued" name=sonobuoy-aggregator namespace=sonobuoy resource=services
I0820 19:48:43.212979      63 handlers.go:72] "Sonobuoy tests were dispatched. Check if there's a CNI is installed." cluster="development-2895"```
```

### Reading the results

To access the cluster it's possible to get the Kube config from the workload with, notice the server IP is internal
so check the Docker port forwarded and change it in the file. After accessing the workload cluster
retrieve the dump logs and access it normally.

```shell
$ kubectl get secrets development-2895-kubeconfig -o json | jq .data.value -r| base64 -d > /tmp/kubeconfig
export KUBECONFIG=/tmp/kubeconfig

# retrieve and uncompress
$ sonobuoy retrieve 
202208201955_sonobuoy_f3b5b48b-c1af-4600-ae74-fd04d0f2d48e.tar.gz

# read the logs
$ less plugins/os-readiness/results/global/out.json 
```

