# ClusterResourceQuota Metrics

| Metric name| Metric type | Labels/tags | Status |
| ---------- | ----------- | ----------- | ----------- |
| openshift_clusterresourcequota_created | Gauge | `name`=&lt;quota-name&gt; | STABLE |
| openshift_clusterresourcequota_labels | Gauge | `name`=&lt;quota-name&gt; | STABLE |
| openshift_clusterresourcequota_selector | Gauge | `name`=&lt;quota-name&gt; <br> `type=`=&lt;annotation\|label&gt; <br> `operator=`=&lt;Operator of MatchExpression, 'In' used for MatchLabels&gt;<br> `key`=&lt;key of annotation or label&gt; <br> `value`=&lt;value&gt; <br>  | STABLE |
| openshift_clusterresourcequota_usage | Gauge | `name`=&lt;quota-name&gt; <br> `resource`=&lt;resource-name&gt; <br> `type`=&lt;hard\|used &gt;| STABLE |
| openshift_clusterresourcequota_ns_usage | Gauge | `name`=&lt;quota-name&gt; <br> `namespace`=&lt;namespace-name&gt; &gt; <br> `resource`=&lt;resource-name&gt; <br> `type`=&lt;hard\|used &gt;| STABLE |