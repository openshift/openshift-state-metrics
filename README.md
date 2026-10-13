# Overview

This repository has been heavily inspired by <https://github.com/kubernetes/kube-state-metrics>.

openshift-state-metrics expands upon kube-state-metrics by adding metrics for OpenShift specific resources.

## How to use

```
$ oc apply -f ./manifests # It will be deployed to openshift-monitoring project
```

## How to generate the manifests

You need make sure jsonnet-bundler and gojsontomal is installed, you can run this make target to install it:

```
$ make $(GOPATH)/bin/jb
$ make $(GOPATH)/bin/gojsontoyaml
```

And then  you can generate the manifests by running:

```
$ make manifests
```

## Documentation

Detailed documentation on the available metrics and usage can be found here: https://github.com/openshift/openshift-state-metrics/blob/master/docs/README.md
