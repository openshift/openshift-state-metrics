# Build Metrics

| Metric name| Metric type | Labels/tags | Status |
| ---------- | ----------- | ----------- | ----------- |
| openshift_build_created | Gauge | `build`=&lt;build-name&gt; <br> `namespace`=&lt;build-namespace&gt; | STABLE |
| openshift_build_metadata_generation | Gauge | `build`=&lt;build-name&gt; <br> `namespace`=&lt;build-namespace&gt; | STABLE |
| openshift_build_labels | Gauge | `build`=&lt;build-name&gt; <br> `namespace`=&lt;build-namespace&gt; | STABLE |
| openshift_build_start | Gauge | `build`=&lt;build-name&gt; <br> `namespace`=&lt;build-namespace&gt; | STABLE |
| openshift_build_complete | Gauge | `build`=&lt;build-name&gt; <br> `namespace`=&lt;build-namespace&gt; | STABLE |
| openshift_build_duration | Gauge | `build`=&lt;build-name&gt; <br> `namespace`=&lt;build-namespace&gt; | STABLE |
| openshift_build_status_phase | Gauge | `build`=&lt;build-name&gt; <br> `namespace`=&lt;build-namespace&gt; <br> `build_phase`=&lt;new\|pending\|running\|error\|failed\|complete\|canceled&gt; |STABLE |