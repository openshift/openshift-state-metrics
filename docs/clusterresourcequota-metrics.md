# ClusterResourceQuota Metrics

| Metric name| Metric type | Labels/tags | Status |
| ---------- | ----------- | ----------- | ----------- |
| openshift_clusterresourcequota_created | Gauge | `name`=&lt;quota-name&gt; | STABLE |
| openshift_clusterresourcequota_labels | Gauge | `name`=&lt;quota-name&gt; | STABLE |
| openshift_clusterresourcequota_selector | Gauge | `name`=&lt;quota-name&gt; <br> `type=`=&lt;annotation\|match-labels\|match-expressions&gt; <br> `operator=`=&lt;Operator only for match-expressions&gt;<br> `key`=&lt;key of annotation or label&gt; <br> `value`=&lt;single value for match-labels and annotations&gt; <br> `values`=&lt;multiple values separated by ',' for match-expressions&gt; <br>  | STABLE |
| openshift_clusterresourcequota_usage | Gauge | `name`=&lt;quota-name&gt; <br> `resource`=&lt;resource-name&gt; <br> `type`=&lt;hard\|used &gt;| STABLE |
| openshift_clusterresourcequota_namespace_usage | Gauge | `name`=&lt;quota-name&gt; <br> `namespace`=&lt;namespace-name&gt; &gt; <br> `resource`=&lt;resource-name&gt; <br> `type`=&lt;hard\|used &gt;| STABLE |