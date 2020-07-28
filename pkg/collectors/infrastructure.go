package collectors

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kube-state-metrics/pkg/metric"
	"k8s.io/kube-state-metrics/pkg/version"

	"github.com/golang/glog"

	configv1 "github.com/openshift/api/config/v1"
	configclient "github.com/openshift/client-go/config/clientset/versioned"
)

var (
	infrastructureMetricFamilies = []metric.FamilyGenerator{
		{
			Name: "cluster_infrastructure_provider",
			Type: metric.MetricTypeGauge,
			Help: "Reports whether the cluster is configured with an infrastructure provider. 'type' is unset if no cloud provider is recognized or set to the constant used by the Infrastructure config. 'region' is set when the cluster clearly identifies a region within the provider. The value is 1 if a cloud provider is set or 0 if it is unset.",
			GenerateFunc: wrapInfrastructureFunc(func(d *configv1.Infrastructure) metric.Family {
				f := metric.Family{}
				var value float64 = 1
				var labelKeys = []string{"type", "region"}
				var labelValues []string

				if status := d.Status.PlatformStatus; status != nil {
					switch {
					// it is illegal to set type to empty string, so let the default case handle
					// empty string (so we can detect it) while preserving the constant None here
					case status.Type == configv1.NonePlatformType:
						labelValues = []string{string(status.Type), ""}
						value = 0
					case status.AWS != nil:
						labelValues = []string{string(status.Type), status.AWS.Region}
					case status.GCP != nil:
						labelValues = []string{string(status.Type), status.GCP.Region}
					default:
						labelValues = []string{string(status.Type), ""}
					}

					f.Metrics = append(f.Metrics, &metric.Metric{
						LabelKeys:   labelKeys,
						LabelValues: labelValues,
						Value:       value,
					})
				}
				return f
			}),
		},
	}
)

func wrapInfrastructureFunc(f func(config *configv1.Infrastructure) metric.Family) func(interface{}) metric.Family {
	return func(obj interface{}) metric.Family {
		return f(obj.(*configv1.Infrastructure))
	}
}

func createInfrastructureListWatch(apiserver string, kubeconfig string) cache.ListWatch {
	configClient, err := createConfigClient(apiserver, kubeconfig)
	if err != nil {
		glog.Fatalf("cannot create config client: %v", err)
	}
	return cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return configClient.ConfigV1().Infrastructures().List(context.TODO(), opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return configClient.ConfigV1().Infrastructures().Watch(context.TODO(), opts)
		},
	}
}

func createConfigClient(apiserver string, kubeconfig string) (*configclient.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags(apiserver, kubeconfig)
	if err != nil {
		return nil, err
	}

	config.UserAgent = version.GetVersion().String()
	config.AcceptContentTypes = "application/vnd.kubernetes.protobuf,application/json"
	config.ContentType = "application/vnd.kubernetes.protobuf"

	client, err := configclient.NewForConfig(config)
	return client, err

}
