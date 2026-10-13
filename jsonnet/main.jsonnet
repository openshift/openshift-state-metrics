local os = (import 'openshift-state-metrics.libsonnet');
{ ['openshift-state-metrics-' + name]: os.openshiftStateMetrics[name] for name in std.objectFields(os.openshiftStateMetrics) }
