# Deployment Metrics

| Metric name| Metric type | Labels/tags | Status |
| ---------- | ----------- | ----------- | ----------- |
| openshift_deploymentconfig_status_replicas | Gauge | `deploymentconfig`=&lt;deploymentconfig-name&gt; <br> `namespace`=&lt;deploymentconfig-namespace&gt; | STABLE |
| openshift_deploymentconfig_status_replicas_available | Gauge | `deploymentconfig`=&lt;deploymentconfig-name&gt; <br> `namespace`=&lt;deploymentconfig-namespace&gt; | STABLE |
| openshift_deploymentconfig_status_replicas_unavailable | Gauge | `deploymentconfig`=&lt;deploymentconfig-name&gt; <br> `namespace`=&lt;deploymentconfig-namespace&gt; | STABLE |
| openshift_deploymentconfig_status_replicas_updated | Gauge | `deploymentconfig`=&lt;deploymentconfig-name&gt; <br> `namespace`=&lt;deploymentconfig-namespace&gt; | STABLE |
| openshift_deploymentconfig_status_observed_generation | Gauge | `deploymentconfig`=&lt;deploymentconfig-name&gt; <br> `namespace`=&lt;deploymentconfig-namespace&gt; | STABLE |
| openshift_deploymentconfig_spec_replicas | Gauge | `deploymentconfig`=&lt;deploymentconfig-name&gt; <br> `namespace`=&lt;deploymentconfig-namespace&gt; | STABLE |
| openshift_deploymentconfig_spec_paused | Gauge | `deploymentconfig`=&lt;deploymentconfig-name&gt; <br> `namespace`=&lt;deploymentconfig-namespace&gt; | STABLE |
| openshift_deploymentconfig_spec_strategy_rollingupdate_max_unavailable | Gauge | `deploymentconfig`=&lt;deploymentconfig-name&gt; <br> `namespace`=&lt;deploymentconfig-namespace&gt; | STABLE |
| openshift_deploymentconfig_spec_strategy_rollingupdate_max_surge | Gauge | `deploymentconfig`=&lt;deploymentconfig-name&gt; <br> `namespace`=&lt;deploymentconfig-namespace&gt; | STABLE |
| openshift_deploymentconfig_metadata_generation | Gauge | `deploymentconfig`=&lt;deploymentconfig-name&gt; <br> `namespace`=&lt;deploymentconfig-namespace&gt; | STABLE |
| openshift_deploymentconfig_labels | Gauge | `deploymentconfig`=&lt;deploymentconfig-name&gt; <br> `namespace`=&lt;deploymentconfig-namespace&gt; | STABLE |
| openshift_deploymentconfig_created | Gauge | `deploymentconfig`=&lt;deploymentconfig-name&gt; <br> `namespace`=&lt;deploymentconfig-namespace&gt; | STABLE |