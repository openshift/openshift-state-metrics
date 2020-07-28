package collectors

import (
	"testing"

	"k8s.io/kube-state-metrics/pkg/metric"

	configv1 "github.com/openshift/api/config/v1"
)

func TestProxyCollector(t *testing.T) {
	// Fixed metadata on type and help text. We prepend this to every expected
	// output so we only have to modify a single place when doing adjustments.
	const metadata = `
		# HELP cluster_proxy_enabled Information about proxy enabled
		# TYPE cluster_proxy_enabled gauge
	`
	cases := []generateMetricsTestCase{
		{
			Obj: &configv1.Proxy{
				Spec: configv1.ProxySpec{
					HTTPProxy: "http-proxy",
				},
			},
			Want: `
        cluster_proxy_enabled{type="http"} 1
		cluster_proxy_enabled{type="https"} 0
		cluster_proxy_enabled{type="trusted_ca"} 0
`,
			MetricNames: []string{"cluster_proxy_enabled"},
		},
	}

	for i, c := range cases {
		c.Func = metric.ComposeMetricGenFuncs(proxyMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
