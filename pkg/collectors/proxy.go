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
	proxyMetricFamilies = []metric.FamilyGenerator{
		{
			Name: "cluster_proxy_enabled",
			Type: metric.MetricTypeGauge,
			Help: "Reports whether the cluster has been configured to use a proxy. 'type' shows the type of proxy configuration that has been set - http for an http proxy, https for an https proxy, and trusted_ca if a custom CA was specified.",
			GenerateFunc: wrapProxyFunc(func(d *configv1.Proxy) metric.Family {
				f := metric.Family{}
				var valueHTTP, valueHTTPS, valueTrustedCA float64
				if len(d.Spec.HTTPProxy) > 0 {
					valueHTTP = 1
				}
				if len(d.Spec.HTTPSProxy) > 0 {
					valueHTTPS = 1
				}
				if len(d.Spec.TrustedCA.Name) > 0 {
					valueTrustedCA = 1
				}

				f.Metrics = append(f.Metrics, &metric.Metric{
					LabelKeys:   []string{"type"},
					LabelValues: []string{"http"},
					Value:       valueHTTP,
				})
				f.Metrics = append(f.Metrics, &metric.Metric{
					LabelKeys:   []string{"type"},
					LabelValues: []string{"https"},
					Value:       valueHTTPS,
				})
				f.Metrics = append(f.Metrics, &metric.Metric{
					LabelKeys:   []string{"type"},
					LabelValues: []string{"trusted_ca"},
					Value:       valueTrustedCA,
				})
				return f
			}),
		},
	}
)

func wrapProxyFunc(f func(config *configv1.Proxy) metric.Family) func(interface{}) metric.Family {
	return func(obj interface{}) metric.Family {
		return f(obj.(*configv1.Proxy))
	}
}

func createProxyListWatch(apiserver string, kubeconfig string) cache.ListWatch {
	configClient, err := createConfigClient(apiserver, kubeconfig)
	if err != nil {
		glog.Fatalf("cannot create config client: %v", err)
	}
	return cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return configClient.ConfigV1().Proxies().List(context.TODO(), opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return configClient.ConfigV1().Proxies().Watch(context.TODO(), opts)
		},
	}
}
