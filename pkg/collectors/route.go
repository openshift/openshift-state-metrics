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

  "github.com/openshift/api/route/v1"
  routeclient "github.com/openshift/client-go/route/clientset/versioned"
)

var (
  descRouteLabelsName          = "openshift_route_labels"
  descRouteLabelsHelp          = "Kubernetes labels converted to Prometheus labels."
  descRouteLabelsDefaultLabels = []string{"namespace", "route_name"}

  routeMetricFamilies = []metric.FamilyGenerator{
    metric.FamilyGenerator{
      Name: "openshift_route_created",
      Type: metric.MetricTypeGauge,
      Help: "Unix creation timestamp",
      GenerateFunc: wrapRouteFunc(func(d *v1.Route) metric.Family {
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
      Name: "openshift_route_info",
      Type: metric.MetricTypeGauge,
      Help: "Information about route.",
      GenerateFunc: wrapRouteFunc(func(d *v1.Route) metric.Family {
        f := metric.Family{}
        var termination string
        for _, m := range d.Status.Ingress {
            if d.Spec.TLS != nil {
              termination = string(d.Spec.TLS.Termination)
            } else {
              termination = ""
            }
          for _, n := range m.Conditions {
            f.Metrics = append(f.Metrics, &metric.Metric{
              LabelKeys: []string{"host_name", "tls_termination", "to_kind", "to_name", "router_name", "conditaion_type", "status"},
              LabelValues: []string{
                d.Spec.Host, 
                termination, 
                d.Spec.To.Kind, 
                d.Spec.To.Name, 
                m.RouterName, 
                string(n.Type), 
                string(n.Status),
              },
              Value: 1,
            })
          }
        }
        
        return f
      }),
    },	
    metric.FamilyGenerator{
      Name: descRouteLabelsName,
      Type: metric.MetricTypeGauge,
      Help: descRouteLabelsHelp,
      GenerateFunc: wrapRouteFunc(func(d *v1.Route) metric.Family {
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

func wrapRouteFunc(f func(*v1.Route) metric.Family) func(interface{}) metric.Family {
  return func(obj interface{}) metric.Family {
    Route := obj.(*v1.Route)

    metricFamily := f(Route)

    for _, m := range metricFamily.Metrics {
      m.LabelKeys = append(descRouteLabelsDefaultLabels, m.LabelKeys...)
      m.LabelValues = append([]string{Route.Namespace, Route.Name}, m.LabelValues...)
    }

    return metricFamily
  }
}

func createRouteListWatch(apiserver string, kubeconfig string, ns string) cache.ListWatch {
  routesclient, err := createRouteClient(apiserver, kubeconfig)
  if err != nil {
    glog.Fatalf("cannot create Route client: %v", err)
  }
  return cache.ListWatch{
    ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
      return routesclient.RouteV1().Routes(ns).List(opts)
    },
    WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
      return routesclient.RouteV1().Routes(ns).Watch(opts)
    },
  }
}

func createRouteClient(apiserver string, kubeconfig string) (*routeclient.Clientset, error) {
  config, err := clientcmd.BuildConfigFromFlags(apiserver, kubeconfig)
  if err != nil {
    return nil, err
  }

  config.UserAgent = version.GetVersion().String()
  config.AcceptContentTypes = "application/vnd.kubernetes.protobuf,application/json"
  config.ContentType = "application/vnd.kubernetes.protobuf"

  client, err := routeclient.NewForConfig(config)
  return client, err

}
