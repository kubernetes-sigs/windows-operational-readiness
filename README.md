Pardon our dust ! This repo is undergoing rapid iteration right now... if you want to run it, 
contact `@jayunit100` on upstream k8s.io slack !

# Operational Readiness Specification for windows

To specify your windows cluster's readiness to run workflows:

- define the example_input.yaml (you can use the example here as a template).
- run the golang program in this repository.

For customization of your specification, see https://github.com/kubernetes/enhancements/pull/2975. 

## First prototype

https://github.com/jayunit100/k8sprototypes/tree/master/windows/op-readiness

## Build and run - development

```
$ make build-test

# run on linux
$ ./op-readiness --provider=local --kubeconfig=<path-to-kubeconfig>

# run on darwin
$ ./op-readiness --os=darwin --provider=local --kubeconfig=<path-to-kubeconfig>
```

## Running tests from a category

Tests categories can be passed in the flag `--category`, this allows users to pick a category of tests by run.
To run ALL tests do not pass the flag.

```
./op-readiness --category core --category networkpolicy

core Ability to access Linux container IPs by service IP (ClusterIP) from Windows containers
...
core Ability to access Windows container IPs by service IP (ClusterIP) from Linux containers
...
core Ability to access Linux container IPs by NodePort IP from Windows containers
```

## Running the Sonobuoy plugin

We support an OCI image and a Sonobuoy plugin, so the user don't need to compile the binary locally
by default the latest version of the E2E binary is builtin the image, if you need to add a custom file
just mount your local version in the plugin at `/app/e2e.test`.

To run the plugin with the default image:

```
make sonobuoy-plugin
```

### Settings a particular category

It's possible to choose one or more categories of tests to run the plugin. In the following
example both `core` and `activedirectory` tests are being use by the plugin.

To allow all tests in the YAML don't use the `--category` flag.

```
spec:
  command:
    - /app/op-readiness
  args:
  - --category 
  - core
  - --category 
  - activedirectory
  - --test-file 
  - tests.yaml
````
