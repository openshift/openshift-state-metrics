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

	"k8s.io/klog/v2"

	v1 "github.com/openshift/api/user/v1"
	groupclient "github.com/openshift/client-go/user/clientset/versioned"
)

var (
	descGroupLabelsName          = "openshift_group_labels"
	descGroupLabelsHelp          = "Kubernetes labels converted to Prometheus labels."
	descGroupLabelsDefaultLabels = []string{"group"}

	groupMetricFamilies = []metric.FamilyGenerator{
		metric.FamilyGenerator{
			Name: "openshift_group_created",
			Type: metric.MetricTypeGauge,
			Help: "Unix creation timestamp",
			GenerateFunc: wrapGroupFunc(func(d *v1.Group) metric.Family {
				f := metric.Family{}

				if !d.CreationTimestamp.IsZero() {
					f.Metrics = append(f.Metrics, &metric.Metric{
						Value: float64(d.CreationTimestamp.Unix()),
					})
				}

				return f
			}),
		},
		metric.FamilyGenerator{
			Name: "openshift_group_user_account",
			Type: metric.MetricTypeGauge,
			Help: "User account in a group.",
			GenerateFunc: wrapGroupFunc(func(d *v1.Group) metric.Family {
				f := metric.Family{}
				if len(d.Users) > 0 {
					for _, user := range d.Users {
						f.Metrics = append(f.Metrics, &metric.Metric{
							LabelKeys:   []string{"user"},
							LabelValues: []string{user},
							Value:       1,
						})
					}

				}
				return f
			}),
		},
	}
)

func wrapGroupFunc(f func(group *v1.Group) metric.Family) func(interface{}) metric.Family {
	return func(obj interface{}) metric.Family {
		group := obj.(*v1.Group)

		metricFamily := f(group)

		for _, m := range metricFamily.Metrics {
			m.LabelKeys = append(descGroupLabelsDefaultLabels, m.LabelKeys...)
			m.LabelValues = append([]string{group.Name}, m.LabelValues...)
		}

		return metricFamily
	}
}

func createGroupListWatch(apiserver string, kubeconfig string, ns string) cache.ListWatch {
	groupclient, err := createGroupClient(apiserver, kubeconfig)
	if err != nil {
		klog.Fatalf("cannot create Group client: %v", err)
	}
	return cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return groupclient.UserV1().Groups().List(context.TODO(), opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return groupclient.UserV1().Groups().Watch(context.TODO(), opts)
		},
	}
}

func createGroupClient(apiserver string, kubeconfig string) (*groupclient.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags(apiserver, kubeconfig)
	if err != nil {
		return nil, err
	}

	config.UserAgent = version.GetVersion().String()
	config.AcceptContentTypes = "application/vnd.kubernetes.protobuf,application/json"
	config.ContentType = "application/vnd.kubernetes.protobuf"

	client, err := groupclient.NewForConfig(config)
	return client, err

}
