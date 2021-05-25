package collectors

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	rcclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"k8s.io/kube-state-metrics/pkg/metric"
	"k8s.io/kube-state-metrics/pkg/version"
)

var (
	descReplicationControllerLabelsDefaultLabels = []string{"namespace", "replicationcontroller"}

	replicationControllerMetricFamilies = []metric.FamilyGenerator{
		metric.FamilyGenerator{
			Name: "openshift_replicationcontroller_owner",
			Type: metric.MetricTypeGauge,
			Help: "Owner of the replication controller",
			GenerateFunc: wrapReplicationControllerFunc(func(d *v1.ReplicationController) metric.Family {
				f := metric.Family{}

				owners := d.GetOwnerReferences()

				if len(owners) == 0 {
					f.Metrics = append(f.Metrics, &metric.Metric{LabelKeys: []string{"owner_kind", "owner_name"},
						LabelValues: []string{"<none>", "<none>"}, Value: 1})
				} else {
					ms := make([]*metric.Metric, len(owners))

					for i, owner := range owners {
						ms[i] = &metric.Metric{LabelValues: []string{owner.Kind, owner.Name},
							LabelKeys: []string{"owner_kind", "owner_name"}, Value: 1}
					}
					f.Metrics = ms
				}
				return f
			}),
		},
	}
)

func wrapReplicationControllerFunc(f func(*v1.ReplicationController) metric.Family) func(interface{}) metric.Family {
	return func(obj interface{}) metric.Family {
		replicationController := obj.(*v1.ReplicationController)

		metricFamily := f(replicationController)

		for _, m := range metricFamily.Metrics {
			m.LabelKeys = append(descReplicationControllerLabelsDefaultLabels, m.LabelKeys...)
			m.LabelValues = append([]string{replicationController.Namespace, replicationController.Name}, m.LabelValues...)
		}

		return metricFamily
	}
}

func createReplicationControllerListWatch(apiserver string, kubeconfig string, ns string) cache.ListWatch {
	replicationControllerClient, err := createReplicationControllerClient(apiserver, kubeconfig)
	if err != nil {
		klog.Fatalf("cannot create Replication Controller client: %v", err)
	}
	return cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return replicationControllerClient.CoreV1().ReplicationControllers(ns).List(context.TODO(), opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return replicationControllerClient.CoreV1().ReplicationControllers(ns).Watch(context.TODO(), opts)
		},
	}
}

func createReplicationControllerClient(apiserver string, kubeconfig string) (*rcclient.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags(apiserver, kubeconfig)
	if err != nil {
		return nil, err
	}

	config.UserAgent = version.GetVersion().String()
	config.AcceptContentTypes = "application/vnd.kubernetes.protobuf,application/json"
	config.ContentType = "application/vnd.kubernetes.protobuf"

	client, err := rcclient.NewForConfig(config)
	return client, err

}
