package collectors

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/api/build/v1"
	"k8s.io/kube-state-metrics/pkg/metric"
)

func TestBuildCollector(t *testing.T) {
	// Fixed metadata on type and help text. We prepend this to every expected
	// output so we only have to modify a single place when doing adjustments.
	const metadata = `
		# HELP openshift_build_created Unix creation timestamp
		# TYPE openshift_build_created gauge
		# HELP openshift_build_metadata_generation Sequence number representing a specific generation of the desired state.
		# TYPE openshift_build_metadata_generation gauge
		# HELP openshift_build_labels Kubernetes labels converted to Prometheus labels.
		# TYPE openshift_build_labels gauge
		# HELP openshift_build_status_phase The build phase
		# TYPE openshift_build_status_phase gauge
		# HELP openshift_build_started Start time of the build
		# TYPE openshift_build_started gauge
		# HELP openshift_build_complete Complete time of the build
		# TYPE openshift_build_complete gauge
		# TYPE openshift_build_duration Duration of the build
`
	cases := []generateMetricsTestCase{
		{
			Obj: &v1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildStatus{
					Phase: v1.BuildPhaseNew,
				},
				Spec: v1.BuildSpec{},
			},
			Want: `
     	openshift_build_complete{build="build1",namespace="ns1"} 0
        openshift_build_created{build="build1",namespace="ns1"} 1.5e+09
        openshift_build_labels{build="build1",label_app="example1",namespace="ns1"} 1
        openshift_build_metadata_generation{build="build1",namespace="ns1"} 21
        openshift_build_start{build="build1",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="cancelled",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="complete",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="error",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="failed",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="new",namespace="ns1"} 1
        openshift_build_status_phase{build="build1",build_phase="pending",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="running",namespace="ns1"} 0
`,

			MetricNames: []string{"openshift_build_created", "openshift_build_metadata_generation", "openshift_build_labels",
				"openshift_build_status_phase", "openshift_build_start", "openshift_build_complete", "openshift_build_duration"},
		},
		{
			Obj: &v1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildStatus{
					Phase: v1.BuildPhasePending,
				},
				Spec: v1.BuildSpec{},
			},
			Want: `
     	openshift_build_complete{build="build1",namespace="ns1"} 0
        openshift_build_created{build="build1",namespace="ns1"} 1.5e+09
        openshift_build_labels{build="build1",label_app="example1",namespace="ns1"} 1
        openshift_build_metadata_generation{build="build1",namespace="ns1"} 21
        openshift_build_start{build="build1",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="cancelled",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="complete",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="error",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="failed",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="new",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="pending",namespace="ns1"} 1
        openshift_build_status_phase{build="build1",build_phase="running",namespace="ns1"} 0
`,
		},
		{
			Obj: &v1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildStatus{
					Phase:          v1.BuildPhaseRunning,
					StartTimestamp: &metav1.Time{Time: time.Unix(1600000000, 0)},
				},
				Spec: v1.BuildSpec{},
			},
			Want: `
		openshift_build_complete{build="build1",namespace="ns1"} 0
        openshift_build_created{build="build1",namespace="ns1"} 1.5e+09
        openshift_build_labels{build="build1",label_app="example1",namespace="ns1"} 1
        openshift_build_metadata_generation{build="build1",namespace="ns1"} 21
        openshift_build_start{build="build1",namespace="ns1"} 1.6e+09
        openshift_build_status_phase{build="build1",build_phase="cancelled",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="complete",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="error",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="failed",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="new",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="pending",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="running",namespace="ns1"} 1
`,
		},
		{
			Obj: &v1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildStatus{
					Phase:          v1.BuildPhaseRunning,
					StartTimestamp: &metav1.Time{Time: time.Unix(1600000000, 0)},
				},
				Spec: v1.BuildSpec{},
			},
			Want: `
		openshift_build_complete{build="build1",namespace="ns1"} 0
        openshift_build_created{build="build1",namespace="ns1"} 1.5e+09
        openshift_build_labels{build="build1",label_app="example1",namespace="ns1"} 1
        openshift_build_metadata_generation{build="build1",namespace="ns1"} 21
        openshift_build_start{build="build1",namespace="ns1"} 1.6e+09
        openshift_build_status_phase{build="build1",build_phase="cancelled",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="complete",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="error",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="failed",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="new",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="pending",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="running",namespace="ns1"} 1
`,
		},
		{
			Obj: &v1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildStatus{
					Phase:               v1.BuildPhaseComplete,
					StartTimestamp:      &metav1.Time{Time: time.Unix(1600000000, 0)},
					CompletionTimestamp: &metav1.Time{Time: time.Unix(1700000000, 0)},
					Duration:            time.Duration(100000000),
				},
				Spec: v1.BuildSpec{},
			},
			Want: `
        openshift_build_complete{build="build1",namespace="ns1"} 1.7e+09
        openshift_build_created{build="build1",namespace="ns1"} 1.5e+09
        openshift_build_duration{build="build1",namespace="ns1"} 1e+08
        openshift_build_labels{build="build1",label_app="example1",namespace="ns1"} 1
        openshift_build_metadata_generation{build="build1",namespace="ns1"} 21
        openshift_build_start{build="build1",namespace="ns1"} 1.6e+09
        openshift_build_status_phase{build="build1",build_phase="cancelled",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="complete",namespace="ns1"} 1
        openshift_build_status_phase{build="build1",build_phase="error",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="failed",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="new",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="pending",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="running",namespace="ns1"} 0
`,
		},
		{
			Obj: &v1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildStatus{
					Phase:               v1.BuildPhaseFailed,
					StartTimestamp:      &metav1.Time{Time: time.Unix(1600000000, 0)},
					CompletionTimestamp: &metav1.Time{Time: time.Unix(1700000000, 0)},
					Duration:            time.Duration(100000000),
				},
				Spec: v1.BuildSpec{},
			},
			Want: `
        openshift_build_complete{build="build1",namespace="ns1"} 1.7e+09
        openshift_build_created{build="build1",namespace="ns1"} 1.5e+09
        openshift_build_duration{build="build1",namespace="ns1"} 1e+08
        openshift_build_labels{build="build1",label_app="example1",namespace="ns1"} 1
        openshift_build_metadata_generation{build="build1",namespace="ns1"} 21
        openshift_build_start{build="build1",namespace="ns1"} 1.6e+09
        openshift_build_status_phase{build="build1",build_phase="cancelled",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="complete",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="error",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="failed",namespace="ns1"} 1
        openshift_build_status_phase{build="build1",build_phase="new",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="pending",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="running",namespace="ns1"} 0
`,
		},
		{
			Obj: &v1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildStatus{
					Phase:               v1.BuildPhaseError,
					StartTimestamp:      &metav1.Time{Time: time.Unix(1600000000, 0)},
					CompletionTimestamp: &metav1.Time{Time: time.Unix(1700000000, 0)},
					Duration:            time.Duration(1000),
				},
				Spec: v1.BuildSpec{},
			},
			Want: `
        openshift_build_complete{build="build1",namespace="ns1"} 1.7e+09
        openshift_build_created{build="build1",namespace="ns1"} 1.5e+09
        openshift_build_duration{build="build1",namespace="ns1"} 1000
        openshift_build_labels{build="build1",label_app="example1",namespace="ns1"} 1
        openshift_build_metadata_generation{build="build1",namespace="ns1"} 21
        openshift_build_start{build="build1",namespace="ns1"} 1.6e+09
        openshift_build_status_phase{build="build1",build_phase="cancelled",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="complete",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="error",namespace="ns1"} 1
        openshift_build_status_phase{build="build1",build_phase="failed",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="new",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="pending",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="running",namespace="ns1"} 0
`,
		},
		{
			Obj: &v1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildStatus{
					Phase:          v1.BuildPhaseCancelled,
					StartTimestamp: &metav1.Time{Time: time.Unix(1600000000, 0)},
				},
				Spec: v1.BuildSpec{},
			},
			Want: `
        openshift_build_complete{build="build1",namespace="ns1"} 0
        openshift_build_created{build="build1",namespace="ns1"} 1.5e+09
        openshift_build_labels{build="build1",label_app="example1",namespace="ns1"} 1
        openshift_build_metadata_generation{build="build1",namespace="ns1"} 21
        openshift_build_start{build="build1",namespace="ns1"} 1.6e+09
        openshift_build_status_phase{build="build1",build_phase="cancelled",namespace="ns1"} 1
        openshift_build_status_phase{build="build1",build_phase="complete",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="error",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="failed",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="new",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="pending",namespace="ns1"} 0
        openshift_build_status_phase{build="build1",build_phase="running",namespace="ns1"} 0
`,
		},
	}

	for i, c := range cases {
		c.Func = metric.ComposeMetricGenFuncs(buildMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
