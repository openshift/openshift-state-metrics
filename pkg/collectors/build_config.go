package collectors

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kube-state-metrics/pkg/metric"
	"k8s.io/kube-state-metrics/pkg/version"

	"github.com/golang/glog"

	"github.com/openshift/api/build/v1"
	buildclient "github.com/openshift/client-go/build/clientset/versioned"
)

var (
	descBuildConfigLabelsName          = "openshift_buildconfig_labels"
	descBuildConfigLabelsHelp          = "Kubernetes labels converted to Prometheus labels."
	descBuildConfigLabelsDefaultLabels = []string{"namespace", "buildconfig"}

	buildconfigMetricFamilies = []metric.FamilyGenerator{
		{
			Name: "openshift_buildconfig_created",
			Type: metric.MetricTypeGauge,
			Help: "Unix creation timestamp",
			GenerateFunc: wrapBuildConfigFunc(func(d *v1.BuildConfig) metric.Family {
				value := float64(0)
				if !d.CreationTimestamp.IsZero() {
					value = float64(d.CreationTimestamp.Unix())
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
			Name: "openshift_buildconfig_metadata_generation",
			Type: metric.MetricTypeGauge,
			Help: "Sequence number representing a specific generation of the desired state.",
			GenerateFunc: wrapBuildConfigFunc(func(d *v1.BuildConfig) metric.Family {
				return metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.ObjectMeta.Generation),
						},
					},
				}
			}),
		},
		{
			Name: descBuildConfigLabelsName,
			Type: metric.MetricTypeGauge,
			Help: descBuildConfigLabelsHelp,
			GenerateFunc: wrapBuildConfigFunc(func(d *v1.BuildConfig) metric.Family {
				labelKeys, labelValues := kubeLabelsToPrometheusLabels(d.Labels)
				return metric.Family{
					Metrics: []*metric.Metric{
						{
							Value:       1,
							LabelKeys:   labelKeys,
							LabelValues: labelValues,
						},
					},
				}
			}),
		},
		metric.FamilyGenerator{
			Name: "openshift_buildconfig_status_latest_version",
			Type: metric.MetricTypeGauge,
			Help: "The latest version of buildconfig.",
			GenerateFunc: wrapBuildConfigFunc(func(d *v1.BuildConfig) metric.Family {
				return metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.Status.LastVersion),
						},
					},
				}
			}),
		},
	}
)

func wrapBuildConfigFunc(f func(config *v1.BuildConfig) metric.Family) func(interface{}) metric.Family {
	return func(obj interface{}) metric.Family {
		buildconfig := obj.(*v1.BuildConfig)

		metricFamily := f(buildconfig)

		for _, m := range metricFamily.Metrics {
			m.LabelKeys = append(descBuildConfigLabelsDefaultLabels, m.LabelKeys...)
			m.LabelValues = append([]string{buildconfig.Namespace, buildconfig.Name}, m.LabelValues...)
		}

		return metricFamily
	}
}

func createBuildConfigListWatch(apiserver string, kubeconfig string, ns string) cache.ListWatch {
	buildclient, err := createBuildClient(apiserver, kubeconfig)
	if err != nil {
		glog.Fatalf("cannot create buildconfig client: %v", err)
	}
	return cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return buildclient.BuildV1().BuildConfigs(ns).List(opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return buildclient.BuildV1().BuildConfigs(ns).Watch(opts)
		},
	}
}

func createBuildClient(apiserver string, kubeconfig string) (*buildclient.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags(apiserver, kubeconfig)
	if err != nil {
		return nil, err
	}

	config.UserAgent = version.GetVersion().String()
	config.AcceptContentTypes = "application/vnd.kubernetes.protobuf,application/json"
	config.ContentType = "application/vnd.kubernetes.protobuf"

	client, err := buildclient.NewForConfig(config)
	return client, err

}
