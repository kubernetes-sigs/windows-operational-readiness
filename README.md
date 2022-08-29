[![Unit Test](https://github.com/kubernetes-sigs/windows-operational-readiness/actions/workflows/unit-tests.yml/badge.svg)](https://github.com/kubernetes-sigs/windows-operational-readiness/actions/workflows/unit-tests.yml)
[![Linter](https://github.com/kubernetes-sigs/windows-operational-readiness/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/kubernetes-sigs/windows-operational-readiness/actions/workflows/golangci-lint.yml)

# Windows Operational Readiness

Define an operational readiness standard for Kubernetes clusters supporting Windows that certifies the readiness of Windows clusters for running production workloads

Related KEP: https://github.com/kubernetes/enhancements/tree/master/keps/sig-windows/2578-windows-conformance

## Build the project

#### Build the project with the default Kubernetes version (e.g. v1.24.0)

```shell
$ make build
```
#### Build the project with a specific Kubernetes version

Compiling the e2e binary from source code:
```shell
$ KUBERNETES_HASH=<Kubernetes commit sha> make build 
```

Or running a pre-compiled released version:

```shell
$ KUBERNETES_VERSION=v1.24.0 make build 
```

## Run the tests

To specify your windows cluster's readiness to run workflows:

- define the tests.yaml (you can use the "tests.yaml" in the repo as a template).

Tests categories can be passed in the flag `--category`, this allows users to pick a category of tests by run.
To run ALL tests do not pass the flag.

```
./op-readiness --provider=local --kubeconfig=<path-to-kubeconfig> --category=Core.Network --category=Sub.NetworkPolicy

Running Operational Readiness Test 1 / 10 : Ability to access Windows container IP by pod IP on Core.Network
...
Running Operational Readiness Test 2 / 10 : Ability to expose windows pods by creating the service ClusterIP on Core.Network
...
```

#### Run the Sonobuoy plugin

We support an OCI image and a Sonobuoy plugin, so the user don't need to compile the binary locally
by default the latest version of the E2E binary is builtin the image, if you need to add a custom file
just mount your local version in the plugin at `/app/e2e.test`.

Before running sonobuoy, taint the Windows worker node. Sonobuoy pod should be scheduled on the control plane node:

```shell
kubectl taint node <windows-worker-node> sonobuoy:NoSchedule
```

To run the plugin with the default image:

```shell
make sonobuoy-plugin
```

To retrieve the sonobuoy result:

```shell
make sonobuoy-results
```

The result can be found in the `./sonobuoy-results` folder.

##### Set a particular category

The `sonobuoy` folder has a [README](sonobuoy/README.md) detailing how to use the templates
to render a custom `sonobuoy-plugin.yaml` file.

#### Running on CAPZ upstream

If you want to test your changes on upstream, use the following bot commmand:

```shell
/test operational-tests-capz-windows-2019
```

## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

You can reach the maintainers of this project at:

- [Slack channel](https://kubernetes.slack.com/messages/sig-windows) 
- [Mailing list](https://groups.google.com/g/kubernetes-sig-windows)

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).

[owners]: https://git.k8s.io/community/contributors/guide/owners.md
[Creative Commons 4.0]: https://git.k8s.io/website/LICENSE
