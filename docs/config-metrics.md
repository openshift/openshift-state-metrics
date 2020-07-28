# Config Metrics

| Metric name| Metric type | Labels/tags | Status |
| ---------- | ----------- | ----------- | ----------- |
| cluster_infrastructure_provider | Gauge | `type`=&lt;cloud-provider-type&gt; <br> `region`=&lt;cloud-provider-region&gt; | STABLE |
| cluster_feature_set | Gauge | `name`=&lt;feature-gate-name&gt; | STABLE |
| cluster_proxy_enabled | Gauge | `type`=&lt;proxy-type&gt; | STABLE |