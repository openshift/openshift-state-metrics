# Documentation

This documentation is intended to be a complete reflection of the current state of the exposed metrics of openshift-state-metrics.

Any contribution to improving this documentation or adding sample usages will be appreciated.

## Table of Contents

- [Metrics Stages](#metrics-stages)
- [Exposed Metrics](#exposed-metrics)
- [CLI arguments](#cli-arguments)

## Metrics Stages

Stages about metrics are grouped into three categories：

| Stage        | Description                                                                                                                |
| ------------ | -------------------------------------------------------------------------------------------------------------------------- |
| EXPERIMENTAL | Metrics which normally correspond to the Kubernetes API object alpha status or spec fields and can be changed at any time. |
| STABLE       | Metrics which should have very few backwards-incompatible changes outside of major version updates.                        |
| DEPRECATED   | Metrics which will be removed once the deprecation timeline is met.                                                        |

## Exposed Metrics

Per group of metrics there is one file for each metrics. See each file for specific documentation about the exposed metrics:

- [BuildConfig Metrics](buildconfig-metrics.md)
- [Build Metrics](build-metrics.md)
- [Config Metrics](config-metrics.md)
- [DeploymentConfig Metrics](deploymentconfig-metrics.md)
- [ClusterResourceQuota Metrics](clusterresourcequota-metrics.md)
- [Route Metrics](route-metrics.md)
- [Group Metrics](group-metrics.md)

## CLI Arguments

Additionally, options for `openshift-state-metrics` can be passed when executing as a CLI, or in a openshift environment. More information can be found here: [CLI Arguments](cli-arguments.md)
