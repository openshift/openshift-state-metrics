package collectors

import (
	"testing"

	"k8s.io/kube-state-metrics/pkg/metric"

	configv1 "github.com/openshift/api/config/v1"
)

func TestFeatureSetCollector(t *testing.T) {
	// Fixed metadata on type and help text. We prepend this to every expected
	// output so we only have to modify a single place when doing adjustments.
	const metadata = `
		# HELP cluster_feature_set Feature Set exposed
		# TYPE cluster_feature_set gauge
	`
	cases := []generateMetricsTestCase{
		{
			Obj: &configv1.FeatureGate{
				Spec: configv1.FeatureGateSpec{
					FeatureGateSelection: configv1.FeatureGateSelection{
						FeatureSet: "feature_set_1",
					},
				},
			},
			Want: `
        cluster_feature_set{name="feature_set_1"} 0	
`,
			MetricNames: []string{"cluster_feature_set"},
		},
	}

	for i, c := range cases {
		c.Func = metric.ComposeMetricGenFuncs(featureSetMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
