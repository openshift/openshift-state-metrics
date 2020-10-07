package collectors

import (
	"context"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kube-state-metrics/pkg/metric"

	"github.com/golang/glog"

	v1 "github.com/openshift/api/build/v1"
)

var (
	descBuildLabelsDefaultLabels = []string{"namespace", "build", "buildconfig", "strategy"}

	buildMetricFamilies = []metric.FamilyGenerator{
		{
			Name: "openshift_build_created_timestamp_seconds",
			Type: metric.MetricTypeGauge,
			Help: "Unix creation timestamp",
			GenerateFunc: wrapBuildFunc(func(b *v1.Build) metric.Family {
				value := float64(0)
				if !b.CreationTimestamp.IsZero() {
					value = float64(b.CreationTimestamp.Unix())
				}
				return metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: value,
						},
					},
				}
			}),
		},
		{
			Name: "openshift_build_metadata_generation_info",
			Type: metric.MetricTypeGauge,
			Help: "Sequence number representing a specific generation of the desired state.",
			GenerateFunc: wrapBuildFunc(func(b *v1.Build) metric.Family {
				return metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(b.ObjectMeta.Generation),
						},
					},
				}
			}),
		},
		{
			Name: "openshift_build_labels",
			Type: metric.MetricTypeGauge,
			Help: "Kubernetes labels converted to Prometheus labels.",
			GenerateFunc: wrapBuildFunc(func(b *v1.Build) metric.Family {
				labelKeys, labelValues := kubeLabelsToPrometheusLabels(b.Labels)
				return metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   labelKeys,
							LabelValues: labelValues,
							Value:       1,
						},
					},
				}
			}),
		},
		{
			Name: "openshift_build_status_phase_total",
			Type: metric.MetricTypeGauge,
			Help: "The build phase.",
			GenerateFunc: wrapBuildFunc(func(b *v1.Build) metric.Family {
				ms := addBuildPhaseMetrics(b.Status.Phase)
				return metric.Family{
					Metrics: ms,
				}
			}),
		},
		{
			Name: "openshift_build_start_timestamp_seconds",
			Type: metric.MetricTypeGauge,
			Help: "Start time of the build",
			GenerateFunc: wrapBuildFunc(func(b *v1.Build) metric.Family {
				value := float64(0)
				if !b.CreationTimestamp.IsZero() && b.Status.StartTimestamp != nil {
					value = float64(b.Status.StartTimestamp.Unix())
				}
				return metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: value,
						},
					},
				}
			}),
		},
		{
			Name: "openshift_build_completed_timestamp_seconds",
			Type: metric.MetricTypeGauge,
			Help: "Completion time of the build",
			GenerateFunc: wrapBuildFunc(func(b *v1.Build) metric.Family {
				value := float64(0)
				if !b.CreationTimestamp.IsZero() && b.Status.CompletionTimestamp != nil {
					value = float64(b.Status.CompletionTimestamp.Unix())
				}
				return metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: value,
						},
					},
				}
			}),
		},
		{
			Name: "openshift_build_duration_seconds",
			Type: metric.MetricTypeGauge,
			Help: "Duration of the build",
			GenerateFunc: wrapBuildFunc(func(b *v1.Build) metric.Family {
				f := metric.Family{}

				if !b.CreationTimestamp.IsZero() && b.Status.Duration != 0 {
					f.Metrics = []*metric.Metric{
						{
							Value: float64(b.Status.Duration / time.Second),
						},
					}
				}
				return f
			}),
		},
	}
)

func wrapBuildFunc(f func(config *v1.Build) metric.Family) func(interface{}) metric.Family {
	return func(obj interface{}) metric.Family {
		build := obj.(*v1.Build)

		metricFamily := f(build)
		buildConfig := determineBuildConfig(build)
		strategy := strings.ToLower(string(build.Spec.Strategy.Type))
		for _, m := range metricFamily.Metrics {
			m.LabelKeys = append(descBuildLabelsDefaultLabels, m.LabelKeys...)
			m.LabelValues = append(
				[]string{build.Namespace, build.Name, buildConfig, strategy},
				m.LabelValues...,
			)
		}

		return metricFamily
	}
}

func determineBuildConfig(build *v1.Build) string {
	if build == nil {
		return ""
	}
	// TODO: Replace quoted strings with build API constants.
	// This requires openshift-state-metrics to be rebased.
	if build.Annotations != nil {
		if _, exists := build.Annotations["openshift.io/build-config.name"]; exists {
			return build.Annotations["openshift.io/build-config.name"]
		}
	}
	if _, exists := build.Labels["openshift.io/build-config.name"]; exists {
		return build.Labels["openshift.io/build-config.name"]
	}
	return build.Labels["buildconfig"]
}

func createBuildListWatch(apiserver string, kubeconfig string, ns string) cache.ListWatch {
	buildclient, err := createBuildClient(apiserver, kubeconfig)
	if err != nil {
		glog.Fatalf("cannot create build client: %v", err)
	}
	return cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return buildclient.BuildV1().Builds(ns).List(context.TODO(), opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return buildclient.BuildV1().Builds(ns).Watch(context.TODO(), opts)
		},
	}
}

// addConditionMetrics generates one metric for each possible node condition
// status. For this function to work properly, the last label in the metric
// description must be the condition.
func addBuildPhaseMetrics(cs v1.BuildPhase) []*metric.Metric {
	return []*metric.Metric{
		{
			LabelValues: []string{"complete"},
			Value:       boolFloat64(cs == v1.BuildPhaseComplete),
			LabelKeys:   []string{"build_phase"},
		},
		{
			LabelValues: []string{"cancelled"},
			Value:       boolFloat64(cs == v1.BuildPhaseCancelled),
			LabelKeys:   []string{"build_phase"},
		},
		{
			LabelValues: []string{"new"},
			Value:       boolFloat64(cs == v1.BuildPhaseNew),
			LabelKeys:   []string{"build_phase"},
		},
		{
			LabelValues: []string{"pending"},
			Value:       boolFloat64(cs == v1.BuildPhasePending),
			LabelKeys:   []string{"build_phase"},
		},
		{
			LabelValues: []string{"running"},
			Value:       boolFloat64(cs == v1.BuildPhaseRunning),
			LabelKeys:   []string{"build_phase"},
		},
		{
			LabelValues: []string{"failed"},
			Value:       boolFloat64(cs == v1.BuildPhaseFailed),
			LabelKeys:   []string{"build_phase"},
		},
		{
			LabelValues: []string{"error"},
			Value:       boolFloat64(cs == v1.BuildPhaseError),
			LabelKeys:   []string{"build_phase"},
		},
	}
}
