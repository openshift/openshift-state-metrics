# ClusterResourceQuota Metrics

| Metric name| Metric type | Labels/tags | Status |
| ---------- | ----------- | ----------- | ----------- |
| openshift_clusterresourcequota_created | Gauge | `name`=&lt;quota-name&gt; | STABLE |
| openshift_clusterresourcequota_labels | Gauge | `name`=&lt;quota-name&gt; | STABLE |
| openshift_clusterresourcequota_usage | Gauge | `name`=&lt;quota-name&gt; <br> `resource`=&lt;resource-name&gt; <br> `type`=&lt;hard\|used &gt;| STABLE |