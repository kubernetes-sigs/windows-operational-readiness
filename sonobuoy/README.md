# Sonobouy Plugin configuration

We use Carvel YTT to render the plugin YAML and patch the required fields. The user must change
the `defaults.yaml` file and rerun `make sonobuoy-config-gen` to render the final `sonobouy-plugin.yaml` file.

## Pre-requisites

1. [Carvel YTT](https://carvel.dev/ytt)

## Changing tests categories

To enable Network and NetworkPolicy categories

```yaml
#@data/values
---
image: winopreadiness/op-readiness:dev
dry-run: false
category:
- Core.Network
- Sub.NetworkPolicy
```