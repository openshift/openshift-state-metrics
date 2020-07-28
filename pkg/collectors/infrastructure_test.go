package collectors

import (
	"testing"

	"k8s.io/kube-state-metrics/pkg/metric"

	configv1 "github.com/openshift/api/config/v1"
)

func TestInfrastructureCollector(t *testing.T) {
	// Fixed metadata on type and help text. We prepend this to every expected
	// output so we only have to modify a single place when doing adjustments.
	const metadata = `
		# HELP cluster_infrastructure_provider Information about cloud infrastructure provider
		# TYPE cluster_infrastructure_provider gauge
	`
	cases := []generateMetricsTestCase{
		{
			Obj: &configv1.Infrastructure{
				Status: configv1.InfrastructureStatus{
					PlatformStatus: &configv1.PlatformStatus{
						Type: "aws",
						AWS: &configv1.AWSPlatformStatus{
							Region: "region_1",
						},
					},
				},
			},
			Want: `
      cluster_infrastructure_provider{region="region_1",type="aws"} 1
`,
			MetricNames: []string{"cluster_infrastructure_provider"},
		},
	}

	for i, c := range cases {
		c.Func = metric.ComposeMetricGenFuncs(infrastructureMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
