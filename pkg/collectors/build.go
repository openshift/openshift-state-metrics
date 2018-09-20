package collectors

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kube-state-metrics/pkg/metric"

	"github.com/golang/glog"

	"github.com/openshift/api/build/v1"
)

var (
	descBuildLabelsName          = "openshift_build_labels"
	descBuildLabelsHelp          = "Kubernetes labels converted to Prometheus labels."
	descBuildLabelsDefaultLabels = []string{"namespace", "build"}

	buildMetricFamilies = []metric.FamilyGenerator{
		metric.FamilyGenerator{
			Name: "openshift_build_created",
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
		metric.FamilyGenerator{
			Name: "openshift_build_metadata_generation",
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
		metric.FamilyGenerator{
			Name: descBuildLabelsName,
			Type: metric.MetricTypeGauge,
			Help: descBuildLabelsHelp,
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
		metric.FamilyGenerator{
			Name: "openshift_build_status_phase",
			Type: metric.MetricTypeGauge,
			Help: "The build phase.",
			GenerateFunc: wrapBuildFunc(func(b *v1.Build) metric.Family {
				ms := addBuildPahseMetrics(b.Status.Phase)
				return metric.Family{
					Metrics: ms,
				}
			}),
		},
		metric.FamilyGenerator{
			Name: "openshift_build_start",
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
		metric.FamilyGenerator{
			Name: "openshift_build_complete",
			Type: metric.MetricTypeGauge,
			Help: "Complet time of the build",
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
		metric.FamilyGenerator{
			Name: "openshift_build_duration",
			Type: metric.MetricTypeGauge,
			Help: "Duration of the build",
			GenerateFunc: wrapBuildFunc(func(b *v1.Build) metric.Family {
				f := metric.Family{}

				if !b.CreationTimestamp.IsZero() && b.Status.Duration != 0 {
					f.Metrics = []*metric.Metric{
						{
							Value: float64(b.Status.Duration),
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

		for _, m := range metricFamily.Metrics {
			m.LabelKeys = append(descBuildLabelsDefaultLabels, m.LabelKeys...)
			m.LabelValues = append([]string{build.Namespace, build.Name}, m.LabelValues...)
		}

		return metricFamily
	}
}

func createBuildListWatch(apiserver string, kubeconfig string, ns string) cache.ListWatch {
	buildclient, err := createBuildClient(apiserver, kubeconfig)
	if err != nil {
		glog.Fatalf("cannot create build client: %v", err)
	}
	return cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return buildclient.BuildV1().Builds(ns).List(opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return buildclient.BuildV1().Builds(ns).Watch(opts)
		},
	}
}

// addConditionMetrics generates one metric for each possible node condition
// status. For this function to work properly, the last label in the metric
// description must be the condition.
func addBuildPahseMetrics(cs v1.BuildPhase) []*metric.Metric {
	return []*metric.Metric{
		&metric.Metric{
			LabelValues: []string{"complete"},
			Value:       boolFloat64(cs == v1.BuildPhaseComplete),
			LabelKeys:   []string{"build_phase"},
		},
		&metric.Metric{
			LabelValues: []string{"cancelled"},
			Value:       boolFloat64(cs == v1.BuildPhaseCancelled),
			LabelKeys:   []string{"build_phase"},
		},
		&metric.Metric{
			LabelValues: []string{"new"},
			Value:       boolFloat64(cs == v1.BuildPhaseNew),
			LabelKeys:   []string{"build_phase"},
		},
		&metric.Metric{
			LabelValues: []string{"pending"},
			Value:       boolFloat64(cs == v1.BuildPhasePending),
			LabelKeys:   []string{"build_phase"},
		},
		&metric.Metric{
			LabelValues: []string{"running"},
			Value:       boolFloat64(cs == v1.BuildPhaseRunning),
			LabelKeys:   []string{"build_phase"},
		},
		&metric.Metric{
			LabelValues: []string{"failed"},
			Value:       boolFloat64(cs == v1.BuildPhaseFailed),
			LabelKeys:   []string{"build_phase"},
		},
		&metric.Metric{
			LabelValues: []string{"error"},
			Value:       boolFloat64(cs == v1.BuildPhaseError),
			LabelKeys:   []string{"build_phase"},
		},
	}
}
