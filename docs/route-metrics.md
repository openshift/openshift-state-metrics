# Route Metrics

| Metric name| Metric type | Labels/tags | Status |
| ---------- | ----------- | ----------- | ----------- |
| openshift_route_created | Gauge | `route`=&lt;route-name&gt; <br> `namespace`=&lt;route-namespace&gt; | STABLE |
| openshift_route_labels | Gauge | `route`=&lt;route-name&gt; <br> `namespace`=&lt;route-namespace&gt; | STABLE |
| openshift_route_info | Gauge | `route`=&lt;route-name&gt; <br> `namespace`=&lt;route-namespace&gt; <br> `host`=&lt;route-host&gt; <br>`path`=&lt;route-path&gt; <br>`tls_termination`=&lt;route-tls-termination&gt; <br> `to_kind`=&lt;route-to-kind&gt; <br>`to-name`=&lt;route-to-name&gt; <br> `to-weight`=&lt;route-to-weight&gt;| STABLE |
| openshift_route_status | Gauge | `route`=&lt;route-name&gt; <br> `namespace`=&lt;route-namespace&gt; <br> `host`=&lt;route-host&gt; <br> `status`=&lt;route-status&gt; <br> `type`=&lt;route-type&gt; <br> `router_name`=&lt;router-name&gt; <br>| STABLE |