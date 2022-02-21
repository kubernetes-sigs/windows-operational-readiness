Pardon our dust ! This repo is undergoing rapid iteration right now... if you want to run it, 
contact jayunit100 on upstream k8s.io slack !

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