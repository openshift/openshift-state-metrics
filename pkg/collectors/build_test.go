package collectors

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/openshift/api/build/v1"
	"k8s.io/kube-state-metrics/pkg/metric"
)

func TestBuildCollector(t *testing.T) {
	// Fixed metadata on type and help text. We prepend this to every expected
	// output so we only have to modify a single place when doing adjustments.
	const metadata = `
		# HELP openshift_build_created_timestamp_seconds Unix creation timestamp
		# TYPE openshift_build_created_timestamp_seconds gauge
		# HELP openshift_build_metadata_generation_info Sequence number representing a specific generation of the desired state.
		# TYPE openshift_build_metadata_generation_info gauge
		# HELP openshift_build_labels Kubernetes labels converted to Prometheus labels.
		# TYPE openshift_build_labels gauge
		# HELP openshift_build_status_phase_total The build phase
		# TYPE openshift_build_status_phase_total gauge
		# HELP openshift_build_start_timestamp_seconds Start time of the build
		# TYPE openshift_build_start_timestamp_seconds gauge
		# HELP openshift_build_completed_timestamp_seconds Complete time of the build
		# TYPE openshift_build_completed_timestamp_seconds gauge
		# TYPE openshift_build_duration_seconds Duration of the build
`
	cases := []generateMetricsTestCase{
		{
			Obj: &v1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Annotations: map[string]string{
						"openshift.io/build-config.name": "build",
					},
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildStatus{
					Phase: v1.BuildPhaseNew,
				},
				Spec: v1.BuildSpec{
					CommonSpec: v1.CommonSpec{
						Strategy: v1.BuildStrategy{
							Type:           v1.DockerBuildStrategyType,
							DockerStrategy: &v1.DockerBuildStrategy{},
						},
					},
				},
			},
			Want: `
     	openshift_build_completed_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_created_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.5e+09
        openshift_build_labels{build="build1",buildconfig="build",label_app="example1",namespace="ns1",strategy="docker"} 1
        openshift_build_metadata_generation_info{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 21
        openshift_build_start_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="cancelled",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="complete",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="error",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="failed",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="new",buildconfig="build",namespace="ns1",strategy="docker"} 1
        openshift_build_status_phase_total{build="build1",build_phase="pending",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="running",buildconfig="build",namespace="ns1",strategy="docker"} 0
`,

			MetricNames: []string{"openshift_build_created_timestamp_seconds", "openshift_build_metadata_generation_info", "openshift_build_labels",
				"openshift_build_status_phase_total", "openshift_build_start_timestamp_seconds", "openshift_build_completed_timestamp_seconds", "openshift_build_duration_seconds"},
		},
		{
			Obj: &v1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Annotations: map[string]string{
						"openshift.io/build-config.name": "build",
					},
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildStatus{
					Phase: v1.BuildPhasePending,
				},
				Spec: v1.BuildSpec{
					CommonSpec: v1.CommonSpec{
						Strategy: v1.BuildStrategy{
							Type:           v1.DockerBuildStrategyType,
							DockerStrategy: &v1.DockerBuildStrategy{},
						},
					},
				},
			},
			Want: `
     	openshift_build_completed_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_created_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.5e+09
        openshift_build_labels{build="build1",buildconfig="build",label_app="example1",namespace="ns1",strategy="docker"} 1
        openshift_build_metadata_generation_info{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 21
        openshift_build_start_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="cancelled",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="complete",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="error",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="failed",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="new",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="pending",buildconfig="build",namespace="ns1",strategy="docker"} 1
        openshift_build_status_phase_total{build="build1",build_phase="running",buildconfig="build",namespace="ns1",strategy="docker"} 0
`,
		},
		{
			Obj: &v1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Annotations: map[string]string{
						"openshift.io/build-config.name": "build",
					},
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildStatus{
					Phase:          v1.BuildPhaseRunning,
					StartTimestamp: &metav1.Time{Time: time.Unix(1600000000, 0)},
				},
				Spec: v1.BuildSpec{
					CommonSpec: v1.CommonSpec{
						Strategy: v1.BuildStrategy{
							Type:           v1.DockerBuildStrategyType,
							DockerStrategy: &v1.DockerBuildStrategy{},
						},
					},
				},
			},
			Want: `
		openshift_build_completed_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_created_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.5e+09
        openshift_build_labels{build="build1",buildconfig="build",label_app="example1",namespace="ns1",strategy="docker"} 1
        openshift_build_metadata_generation_info{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 21
        openshift_build_start_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.6e+09
        openshift_build_status_phase_total{build="build1",build_phase="cancelled",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="complete",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="error",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="failed",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="new",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="pending",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="running",buildconfig="build",namespace="ns1",strategy="docker"} 1
`,
		},
		{
			Obj: &v1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Annotations: map[string]string{
						"openshift.io/build-config.name": "build",
					},
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildStatus{
					Phase:          v1.BuildPhaseRunning,
					StartTimestamp: &metav1.Time{Time: time.Unix(1600000000, 0)},
				},
				Spec: v1.BuildSpec{
					CommonSpec: v1.CommonSpec{
						Strategy: v1.BuildStrategy{
							Type:           v1.DockerBuildStrategyType,
							DockerStrategy: &v1.DockerBuildStrategy{},
						},
					},
				},
			},
			Want: `
		openshift_build_completed_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_created_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.5e+09
        openshift_build_labels{build="build1",buildconfig="build",label_app="example1",namespace="ns1",strategy="docker"} 1
        openshift_build_metadata_generation_info{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 21
        openshift_build_start_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.6e+09
        openshift_build_status_phase_total{build="build1",build_phase="cancelled",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="complete",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="error",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="failed",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="new",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="pending",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="running",buildconfig="build",namespace="ns1",strategy="docker"} 1
`,
		},
		{
			Obj: &v1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Annotations: map[string]string{
						"openshift.io/build-config.name": "build",
					},
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildStatus{
					Phase:               v1.BuildPhaseComplete,
					StartTimestamp:      &metav1.Time{Time: time.Unix(1600000000, 0)},
					CompletionTimestamp: &metav1.Time{Time: time.Unix(1700000000, 0)},
					Duration:            10 * time.Second,
				},
				Spec: v1.BuildSpec{
					CommonSpec: v1.CommonSpec{
						Strategy: v1.BuildStrategy{
							Type:           v1.DockerBuildStrategyType,
							DockerStrategy: &v1.DockerBuildStrategy{},
						},
					},
				},
			},
			Want: `
        openshift_build_completed_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.7e+09
        openshift_build_created_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.5e+09
        openshift_build_duration_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 10
        openshift_build_labels{build="build1",buildconfig="build",label_app="example1",namespace="ns1",strategy="docker"} 1
        openshift_build_metadata_generation_info{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 21
        openshift_build_start_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.6e+09
        openshift_build_status_phase_total{build="build1",build_phase="cancelled",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="complete",buildconfig="build",namespace="ns1",strategy="docker"} 1
        openshift_build_status_phase_total{build="build1",build_phase="error",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="failed",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="new",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="pending",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="running",buildconfig="build",namespace="ns1",strategy="docker"} 0
`,
		},
		{
			Obj: &v1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Annotations: map[string]string{
						"openshift.io/build-config.name": "build",
					},
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildStatus{
					Phase:               v1.BuildPhaseFailed,
					StartTimestamp:      &metav1.Time{Time: time.Unix(1600000000, 0)},
					CompletionTimestamp: &metav1.Time{Time: time.Unix(1700000000, 0)},
					Duration:            10 * time.Second,
				},
				Spec: v1.BuildSpec{
					CommonSpec: v1.CommonSpec{
						Strategy: v1.BuildStrategy{
							Type:           v1.DockerBuildStrategyType,
							DockerStrategy: &v1.DockerBuildStrategy{},
						},
					},
				},
			},
			Want: `
        openshift_build_completed_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.7e+09
        openshift_build_created_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.5e+09
        openshift_build_duration_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 10
        openshift_build_labels{build="build1",buildconfig="build",label_app="example1",namespace="ns1",strategy="docker"} 1
        openshift_build_metadata_generation_info{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 21
        openshift_build_start_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.6e+09
        openshift_build_status_phase_total{build="build1",build_phase="cancelled",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="complete",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="error",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="failed",buildconfig="build",namespace="ns1",strategy="docker"} 1
        openshift_build_status_phase_total{build="build1",build_phase="new",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="pending",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="running",buildconfig="build",namespace="ns1",strategy="docker"} 0
`,
		},
		{
			Obj: &v1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Annotations: map[string]string{
						"openshift.io/build-config.name": "build",
					},
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildStatus{
					Phase:               v1.BuildPhaseError,
					StartTimestamp:      &metav1.Time{Time: time.Unix(1600000000, 0)},
					CompletionTimestamp: &metav1.Time{Time: time.Unix(1700000000, 0)},
					Duration:            1 * time.Second,
				},
				Spec: v1.BuildSpec{
					CommonSpec: v1.CommonSpec{
						Strategy: v1.BuildStrategy{
							Type:           v1.DockerBuildStrategyType,
							DockerStrategy: &v1.DockerBuildStrategy{},
						},
					},
				},
			},
			Want: `
        openshift_build_completed_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.7e+09
        openshift_build_created_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.5e+09
        openshift_build_duration_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1
        openshift_build_labels{build="build1",buildconfig="build",label_app="example1",namespace="ns1",strategy="docker"} 1
        openshift_build_metadata_generation_info{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 21
        openshift_build_start_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.6e+09
        openshift_build_status_phase_total{build="build1",build_phase="cancelled",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="complete",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="error",buildconfig="build",namespace="ns1",strategy="docker"} 1
        openshift_build_status_phase_total{build="build1",build_phase="failed",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="new",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="pending",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="running",buildconfig="build",namespace="ns1",strategy="docker"} 0
`,
		},
		{
			Obj: &v1.Build{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "build1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Annotations: map[string]string{
						"openshift.io/build-config.name": "build",
					},
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.BuildStatus{
					Phase:          v1.BuildPhaseCancelled,
					StartTimestamp: &metav1.Time{Time: time.Unix(1600000000, 0)},
				},
				Spec: v1.BuildSpec{
					CommonSpec: v1.CommonSpec{
						Strategy: v1.BuildStrategy{
							Type:           v1.DockerBuildStrategyType,
							DockerStrategy: &v1.DockerBuildStrategy{},
						},
					},
				},
			},
			Want: `
        openshift_build_completed_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_created_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.5e+09
        openshift_build_labels{build="build1",buildconfig="build",label_app="example1",namespace="ns1",strategy="docker"} 1
        openshift_build_metadata_generation_info{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 21
        openshift_build_start_timestamp_seconds{build="build1",buildconfig="build",namespace="ns1",strategy="docker"} 1.6e+09
        openshift_build_status_phase_total{build="build1",build_phase="cancelled",buildconfig="build",namespace="ns1",strategy="docker"} 1
        openshift_build_status_phase_total{build="build1",build_phase="complete",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="error",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="failed",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="new",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="pending",buildconfig="build",namespace="ns1",strategy="docker"} 0
        openshift_build_status_phase_total{build="build1",build_phase="running",buildconfig="build",namespace="ns1",strategy="docker"} 0
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
