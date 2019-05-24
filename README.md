# Overview

This repository has been heavily inspired by <https://github.com/kubernetes/kube-state-metrics>.

openshift-state-metrics expands upon kube-state-metrics by adding metrics for OpenShift specific resources.

## How to use

```
$ oc apply -f ./manifests # It will be deployed to openshift-monitoring project
```

## Documentation

Detailed documentation on the available metrics and usage can be found here: https://github.com/openshift/openshift-state-metrics/blob/master/docs/README.md
