# Build Metrics

| Metric name| Metric type | Labels/tags | Status |
| ---------- | ----------- | ----------- | ----------- |
| openshift_build_created_timestamp_seconds | Gauge | `build`=&lt;build-name&gt; <br> `buildconfig`=&lt;build-config&gt; <br> `namespace`=&lt;build-namespace&gt; <br> `strategy`=&lt;custom\|docker\|jenkinspipeline\|source&gt; | STABLE |
| openshift_build_metadata_generation_info | Gauge | `build`=&lt;build-name&gt; <br> `buildconfig`=&lt;build-config&gt; <br> `namespace`=&lt;build-namespace&gt; <br> `strategy`=&lt;custom\|docker\|jenkinspipeline\|source&gt; | STABLE |
| openshift_build_labels | Gauge | `build`=&lt;build-name&gt; <br> `buildconfig`=&lt;build-config&gt; <br> `namespace`=&lt;build-namespace&gt; <br> `strategy`=&lt;custom\|docker\|jenkinspipeline\|source&gt; | STABLE |
| openshift_build_start_timestamp_seconds | Gauge | `build`=&lt;build-name&gt; <br> `buildconfig`=&lt;build-config&gt; <br> `namespace`=&lt;build-namespace&gt; <br> `strategy`=&lt;custom\|docker\|jenkinspipeline\|source&gt; | STABLE |
| openshift_build_completed_timestamp_seconds | Gauge | `build`=&lt;build-name&gt; <br> `buildconfig`=&lt;build-config&gt; <br> `namespace`=&lt;build-namespace&gt; <br> `strategy`=&lt;custom\|docker\|jenkinspipeline\|source&gt; | STABLE |
| openshift_build_duration_seconds | Gauge | `build`=&lt;build-name&gt; <br> `buildconfig`=&lt;build-config&gt; <br> `namespace`=&lt;build-namespace&gt; <br> `strategy`=&lt;custom\|docker\|jenkinspipeline\|source&gt; | STABLE |
| openshift_build_status_phase_total | Gauge | `build`=&lt;build-name&gt; <br> `build_phase`=&lt;new\|pending\|running\|error\|failed\|complete\|canceled&gt; <br> `buildconfig`=&lt;build-config&gt; <br> `namespace`=&lt;build-namespace&gt; <br> `strategy`=&lt;custom\|docker\|jenkinspipeline\|source&gt; |STABLE |