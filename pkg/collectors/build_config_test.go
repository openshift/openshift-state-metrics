package collectors

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/api/build/v1"
	"k8s.io/kube-state-metrics/pkg/metric"
)

func TestBuildConfigCollector(t *testing.T) {
	// Fixed metadata on type and help text. We prepend this to every expected
	// output so we only have to modify a single place when doing adjustments.
	const metadata = `
		# HELP openshift_buildconfig_created Unix creation timestamp
		# TYPE openshift_buildconfig_created gauge
		# HELP openshift_buildconfig_metadata_generation Sequence number representing a specific generation of the desired state.
		# TYPE openshift_buildconfig_metadata_generation gauge
		# HELP openshift_buildconfig_labels Kubernetes labels converted to Prometheus labels.
		# TYPE openshift_buildconfig_labels gauge
		# HELP openshift_buildconfig_status_latest_version The latest version of buildconfig.
		# TYPE openshift_buildconfig_status_latest_version gauge
	`
	cases := []generateMetricsTestCase{
		{
			Obj: &v1.BuildConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildConfigStatus{
					LastVersion: int64(1),
				},
				Spec: v1.BuildConfigSpec{},
			},
			Want: `
        openshift_buildconfig_created{buildconfig="build1",namespace="ns1"} 1.5e+09
        openshift_buildconfig_labels{buildconfig="build1",label_app="example1",namespace="ns1"} 1
        openshift_buildconfig_metadata_generation{buildconfig="build1",namespace="ns1"} 21
        openshift_buildconfig_status_latest_version{buildconfig="build1",namespace="ns1"} 1
`,
			MetricNames: []string{"openshift_buildconfig_labels", "openshift_buildconfig_status_latest_version", "openshift_buildconfig_created", "openshift_buildconfig_metadata_generation"},
		},
	}

	for i, c := range cases {
		c.Func = metric.ComposeMetricGenFuncs(buildconfigMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
