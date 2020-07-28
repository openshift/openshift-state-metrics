package collectors

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/kube-state-metrics/pkg/metric"

	"github.com/golang/glog"

	configv1 "github.com/openshift/api/config/v1"
)

var (
	featureSetMetricFamilies = []metric.FamilyGenerator{
		{
			Name: "cluster_feature_set",
			Type: metric.MetricTypeGauge,
			Help: "Reports the feature set the cluster is configured to expose. 'name' corresponds to the name of the feature gate",
			GenerateFunc: wrapFeatureGateFunc(func(d *configv1.FeatureGate) metric.Family {
				f := metric.Family{}
				var value float64
				if d.Spec.FeatureSet == configv1.Default {
					value = 1
				}
				f.Metrics = append(f.Metrics, &metric.Metric{
					LabelKeys:   []string{"name"},
					LabelValues: []string{string(d.Spec.FeatureSet)},
					Value:       value,
				})
				return f
			}),
		},
	}
)

func wrapFeatureGateFunc(f func(config *configv1.FeatureGate) metric.Family) func(interface{}) metric.Family {
	return func(obj interface{}) metric.Family {
		return f(obj.(*configv1.FeatureGate))
	}
}

func createFeatureGateListWatch(apiserver string, kubeconfig string) cache.ListWatch {
	configClient, err := createConfigClient(apiserver, kubeconfig)
	if err != nil {
		glog.Fatalf("cannot create config client: %v", err)
	}
	return cache.ListWatch{

		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return configClient.ConfigV1().FeatureGates().List(context.TODO(), opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return configClient.ConfigV1().FeatureGates().Watch(context.TODO(), opts)
		},
	}
}
