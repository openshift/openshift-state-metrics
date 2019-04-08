package collectors

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/kube-state-metrics/pkg/metric"
	"k8s.io/kube-state-metrics/pkg/version"

	"github.com/golang/glog"

	"github.com/openshift/api/apps/v1"
	appsclient "github.com/openshift/client-go/apps/clientset/versioned"
)

var (
	descDeploymentLabelsName          = "openshift_deploymentconfig_labels"
	descDeploymentLabelsHelp          = "Kubernetes labels converted to Prometheus labels."
	descDeploymentLabelsDefaultLabels = []string{"namespace", "deploymentconfig"}

	deploymentMetricFamilies = []metric.FamilyGenerator{
		metric.FamilyGenerator{
			Name: "openshift_deploymentconfig_created",
			Type: metric.MetricTypeGauge,
			Help: "Unix creation timestamp",
			GenerateFunc: wrapDeploymentFunc(func(d *v1.DeploymentConfig) metric.Family {
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
			Name: "openshift_deploymentconfig_status_replicas",
			Type: metric.MetricTypeGauge,
			Help: "The number of replicas per deployment.",
			GenerateFunc: wrapDeploymentFunc(func(d *v1.DeploymentConfig) metric.Family {
				return metric.Family{
					Metrics: []*metric.Metric{
						{
							Value: float64(d.Status.Replicas),
						},
					}}
			}),
		},
		metric.FamilyGenerator{
			Name: "openshift_deploymentconfig_status_replicas_available",
			Type: metric.MetricTypeGauge,
			Help: "The number of available replicas per deployment.",
			GenerateFunc: wrapDeploymentFunc(func(d *v1.DeploymentConfig) metric.Family {
				return metric.Family{Metrics: []*metric.Metric{
					{
						Value: float64(d.Status.AvailableReplicas),
					},
				}}

			}),
		},
		metric.FamilyGenerator{
			Name: "openshift_deploymentconfig_status_replicas_unavailable",
			Type: metric.MetricTypeGauge,
			Help: "The number of unavailable replicas per deployment.",
			GenerateFunc: wrapDeploymentFunc(func(d *v1.DeploymentConfig) metric.Family {
				return metric.Family{Metrics: []*metric.Metric{
					{
						Value: float64(d.Status.UnavailableReplicas),
					},
				}}
			}),
		},
		metric.FamilyGenerator{
			Name: "openshift_deploymentconfig_status_replicas_updated",
			Type: metric.MetricTypeGauge,
			Help: "The number of updated replicas per deployment.",
			GenerateFunc: wrapDeploymentFunc(func(d *v1.DeploymentConfig) metric.Family {
				return metric.Family{Metrics: []*metric.Metric{
					{
						Value: float64(d.Status.UpdatedReplicas),
					},
				}}
			}),
		},
		metric.FamilyGenerator{
			Name: "openshift_deploymentconfig_status_observed_generation",
			Type: metric.MetricTypeGauge,
			Help: "The generation observed by the deployment controller.",
			GenerateFunc: wrapDeploymentFunc(func(d *v1.DeploymentConfig) metric.Family {
				return metric.Family{Metrics: []*metric.Metric{
					{
						Value: float64(d.Status.ObservedGeneration),
					},
				}}
			}),
		},
		metric.FamilyGenerator{
			Name: "openshift_deploymentconfig_spec_replicas",
			Type: metric.MetricTypeGauge,
			Help: "Number of desired pods for a deployment.",
			GenerateFunc: wrapDeploymentFunc(func(d *v1.DeploymentConfig) metric.Family {
				return metric.Family{Metrics: []*metric.Metric{
					{
						Value: float64(d.Spec.Replicas),
					},
				}}
			}),
		},
		metric.FamilyGenerator{
			Name: "openshift_deploymentconfig_spec_paused",
			Type: metric.MetricTypeGauge,
			Help: "Whether the deployment is paused and will not be processed by the deployment controller.",
			GenerateFunc: wrapDeploymentFunc(func(d *v1.DeploymentConfig) metric.Family {
				return metric.Family{Metrics: []*metric.Metric{
					{
						Value: boolFloat64(d.Spec.Paused),
					},
				}}
			}),
		},
		metric.FamilyGenerator{
			Name: "openshift_deploymentconfig_spec_strategy_rollingupdate_max_unavailable",
			Type: metric.MetricTypeGauge,
			Help: "Maximum number of unavailable replicas during a rolling update of a deployment.",
			GenerateFunc: wrapDeploymentFunc(func(d *v1.DeploymentConfig) metric.Family {
				if d.Spec.Strategy.RollingParams == nil {
					return metric.Family{}
				}

				maxUnavailable, err := intstr.GetValueFromIntOrPercent(d.Spec.Strategy.RollingParams.MaxUnavailable, int(d.Spec.Replicas), true)
				if err != nil {
					panic(err)
				}
				return metric.Family{Metrics: []*metric.Metric{
					{
						Value: float64(maxUnavailable),
					},
				}}
			}),
		},
		metric.FamilyGenerator{
			Name: "openshift_deploymentconfig_spec_strategy_rollingupdate_max_surge",
			Type: metric.MetricTypeGauge,
			Help: "Maximum number of replicas that can be scheduled above the desired number of replicas during a rolling update of a deployment.",
			GenerateFunc: wrapDeploymentFunc(func(d *v1.DeploymentConfig) metric.Family {
				if d.Spec.Strategy.RollingParams == nil {
					return metric.Family{}
				}

				maxSurge, err := intstr.GetValueFromIntOrPercent(d.Spec.Strategy.RollingParams.MaxSurge, int(d.Spec.Replicas), true)
				if err != nil {
					panic(err)
				}
				return metric.Family{Metrics: []*metric.Metric{
					{
						Value: float64(maxSurge),
					},
				}}
			}),
		},
		metric.FamilyGenerator{
			Name: "openshift_deploymentconfig_metadata_generation",
			Type: metric.MetricTypeGauge,
			Help: "Sequence number representing a specific generation of the desired state.",
			GenerateFunc: wrapDeploymentFunc(func(d *v1.DeploymentConfig) metric.Family {
				return metric.Family{Metrics: []*metric.Metric{
					{
						Value: float64(d.ObjectMeta.Generation),
					},
				}}
			}),
		},
		metric.FamilyGenerator{
			Name: descDeploymentLabelsName,
			Type: metric.MetricTypeGauge,
			Help: descDeploymentLabelsHelp,
			GenerateFunc: wrapDeploymentFunc(func(d *v1.DeploymentConfig) metric.Family {
				labelKeys, labelValues := kubeLabelsToPrometheusLabels(d.Labels)
				return metric.Family{Metrics: []*metric.Metric{
					{
						LabelKeys:   labelKeys,
						LabelValues: labelValues,
						Value:       1,
					},
				}}
			}),
		},
	}
)

func wrapDeploymentFunc(f func(*v1.DeploymentConfig) metric.Family) func(interface{}) metric.Family {
	return func(obj interface{}) metric.Family {
		deployment := obj.(*v1.DeploymentConfig)

		metricFamily := f(deployment)

		for _, m := range metricFamily.Metrics {
			m.LabelKeys = append(descDeploymentLabelsDefaultLabels, m.LabelKeys...)
			m.LabelValues = append([]string{deployment.Namespace, deployment.Name}, m.LabelValues...)
		}

		return metricFamily
	}
}

func createDeploymentListWatch(apiserver string, kubeconfig string, ns string) cache.ListWatch {
	appsclient, err := createAppsClient(apiserver, kubeconfig)
	if err != nil {
		glog.Fatalf("cannot create deploymentconfig client: %v", err)
	}
	return cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return appsclient.AppsV1().DeploymentConfigs(ns).List(opts)
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return appsclient.AppsV1().DeploymentConfigs(ns).Watch(opts)
		},
	}
}

func createAppsClient(apiserver string, kubeconfig string) (*appsclient.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags(apiserver, kubeconfig)
	if err != nil {
		return nil, err
	}

	config.UserAgent = version.GetVersion().String()
	config.AcceptContentTypes = "application/vnd.kubernetes.protobuf,application/json"
	config.ContentType = "application/vnd.kubernetes.protobuf"

	client, err := appsclient.NewForConfig(config)
	return client, err

}
