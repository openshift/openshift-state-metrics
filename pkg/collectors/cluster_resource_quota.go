package collectors

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kube-state-metrics/pkg/metric"
	"strings"

	"github.com/golang/glog"
	"github.com/openshift/api/quota/v1"
	quotaclient "github.com/openshift/client-go/quota/clientset/versioned"
	"k8s.io/kube-state-metrics/pkg/version"
)

var (
	descClusterResourceQuotaLabelsName          = "openshift_clusterresourcequota_labels"
	descClusterResourceQuotaLabelsHelp          = "Kubernetes labels converted to Prometheus labels."
	descClusterResourceQuotaLabelsDefaultLabels = []string{"name"}

	quotaMetricFamilies = []metric.FamilyGenerator{
		metric.FamilyGenerator{
			Name: "openshift_clusterresourcequota_created",
			Type: metric.MetricTypeGauge,
			Help: "Unix creation timestamp",
			GenerateFunc: wrapClusterResourceQuotaFunc(func(b *v1.ClusterResourceQuota) metric.Family {
				f := metric.Family{}
				if !b.CreationTimestamp.IsZero() {
					f.Metrics = append(f.Metrics, &metric.Metric{
						Value: float64(b.CreationTimestamp.Unix()),
					})
				}
				return f
			}),
		},
		metric.FamilyGenerator{
			Name: descClusterResourceQuotaLabelsName,
			Type: metric.MetricTypeGauge,
			Help: descClusterResourceQuotaLabelsHelp,
			GenerateFunc: wrapClusterResourceQuotaFunc(func(quota *v1.ClusterResourceQuota) metric.Family {
				labelKeys, labelValues := kubeLabelsToPrometheusLabels(quota.Labels)
				return metric.Family{
					Metrics: []*metric.Metric{
						{
							LabelKeys:   labelKeys,
							LabelValues: labelValues,
							Value:       1,
						},
					}}
			}),
		},
		metric.FamilyGenerator{
			Name: "openshift_clusterresourcequota_usage",
			Type: metric.MetricTypeGauge,
			Help: "Usage about resource quota.",
			GenerateFunc: wrapClusterResourceQuotaFunc(func(r *v1.ClusterResourceQuota) metric.Family {
				f := metric.Family{}

				for res, qty := range r.Spec.Quota.Hard {
					f.Metrics = append(f.Metrics, &metric.Metric{
						LabelValues: []string{string(res), "hard"},
						Value:       float64(qty.MilliValue()) / 1000,
					})
				}
				for res, qty := range r.Status.Total.Used {
					f.Metrics = append(f.Metrics, &metric.Metric{
						LabelValues: []string{string(res), "used"},
						Value:       float64(qty.MilliValue()) / 1000,
					})
				}

				for _, m := range f.Metrics {
					m.LabelKeys = []string{"resource", "type"}
				}

				return f
			}),
		},
		metric.FamilyGenerator{
			Name: "openshift_clusterresourcequota_namespace_usage",
			Type: metric.MetricTypeGauge,
			Help: "Usage about clusterresource quota per namespace.",
			GenerateFunc: wrapClusterResourceQuotaFunc(func(r *v1.ClusterResourceQuota) metric.Family {
				f := metric.Family{}

				for _, rq := range r.Status.Namespaces {
					for res, qty := range rq.Status.Hard {
						f.Metrics = append(f.Metrics, &metric.Metric{
							LabelValues: []string{string(rq.Namespace), string(res), "hard"},
							Value:       float64(qty.MilliValue()) / 1000,
						})
					}
					for res, qty := range rq.Status.Used {
						f.Metrics = append(f.Metrics, &metric.Metric{
							LabelValues: []string{string(rq.Namespace), string(res), "used"},
							Value:       float64(qty.MilliValue()) / 1000,
						})
					}

					for _, m := range f.Metrics {
						m.LabelKeys = []string{"namespace", "resource", "type"}
					}
				}

				return f
			}),
		},
		metric.FamilyGenerator{
			Name: "openshift_clusterresourcequota_selector",
			Type: metric.MetricTypeGauge,
			Help: "Selector of clusterresource quota, which defines the affected namespaces.",
			GenerateFunc: wrapClusterResourceQuotaFunc(func(r *v1.ClusterResourceQuota) metric.Family {
				f := metric.Family{}

				sel := r.Spec.Selector
				labelKeys := []string{"type", "key", "values"}
				for key, value := range sel.AnnotationSelector {
					f.Metrics = append(f.Metrics, &metric.Metric{
						LabelKeys:   labelKeys,
						LabelValues: []string{"annotation", string(key), string(value)},
						Value:       float64(1),
					})
				}

				if sel.LabelSelector != nil {
					labelKeys := []string{"type", "key", "value"}

					for key, value := range sel.LabelSelector.MatchLabels {
						f.Metrics = append(f.Metrics, &metric.Metric{
							LabelKeys:   labelKeys,
							LabelValues: []string{"match-labels", string(key), string(value)},
							Value:       float64(1),
						})
					}

					labelKeys = []string{"type", "operator", "key", "values"}
					for _, labelReq := range sel.LabelSelector.MatchExpressions {
						f.Metrics = append(f.Metrics, &metric.Metric{
							LabelKeys:   labelKeys,
							LabelValues: []string{"match-expressions", string(labelReq.Operator), string(labelReq.Key), string(strings.Join(labelReq.Values, ","))},
							Value:       float64(1),
						})

					}

				}

				return f
			}),
		},
	}
)

func wrapClusterResourceQuotaFunc(f func(config *v1.ClusterResourceQuota) metric.Family) func(interface{}) metric.Family {
	return func(obj interface{}) metric.Family {
		quota := obj.(*v1.ClusterResourceQuota)

		metricFamily := f(quota)

		for _, m := range metricFamily.Metrics {
			m.LabelKeys = append(descClusterResourceQuotaLabelsDefaultLabels, m.LabelKeys...)
			m.LabelValues = append([]string{quota.Name}, m.LabelValues...)
		}

		return metricFamily
	}
}

func createClusterResourceQuotaListWatch(apiserver string, kubeconfig string, ns string) cache.ListWatch {
	quotaclient, err := createClusterResourceQuotaClient(apiserver, kubeconfig)
	if err != nil {
		glog.Fatalf("cannot create quota client: %v", err)
	}
	return cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return quotaclient.QuotaV1().ClusterResourceQuotas().List(opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return quotaclient.QuotaV1().ClusterResourceQuotas().Watch(opts)
		},
	}
}

func createClusterResourceQuotaClient(apiserver string, kubeconfig string) (*quotaclient.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags(apiserver, kubeconfig)
	if err != nil {
		return nil, err
	}

	config.UserAgent = version.GetVersion().String()
	config.AcceptContentTypes = "application/vnd.kubernetes.protobuf,application/json"
	config.ContentType = "application/vnd.kubernetes.protobuf"

	client, err := quotaclient.NewForConfig(config)
	return client, err

}
