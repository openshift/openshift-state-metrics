package collectors

import (
	"testing"
	"time"

	"github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kube-state-metrics/pkg/metric"
)

var (
	weight int32 = 100
)

func TestRouteConfigCollector(t *testing.T) {
	// Fixed metadata on type and help text. We prepend this to every expected
	// output so we only have to modify a single place when doing adjustments.
	const metadata = `
		# HELP openshift_route_created Unix creation timestamp
		# TYPE openshift_route_created gauge
		# HELP openshift_route_labels Kubernetes labels converted to Prometheus labels.
		# TYPE openshift_route_labels gauge
		# HELP openshift_route_info Information about route.
		# TYPE openshift_route_info gauge
	`
	cases := []generateMetricsTestCase{
		{
			Obj: &v1.Route{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "route1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Labels: map[string]string{
						"app": "good1",
					},
				},
				Status: v1.RouteStatus{
					Ingress: []v1.RouteIngress{
						{
							Host:       "example.com",
							RouterName: "router1",
							Conditions: []v1.RouteIngressCondition{
								{
									Type:   v1.RouteAdmitted,
									Status: corev1.ConditionTrue,
								},
							},
						},
						{
							Host:       "example.com",
							RouterName: "router2",
							Conditions: []v1.RouteIngressCondition{
								{
									Type:   v1.RouteAdmitted,
									Status: corev1.ConditionTrue,
								},
							},
						},
					},
				},
				Spec: v1.RouteSpec{
					Host: "example.com",
					TLS: &v1.TLSConfig{
						Termination: "edge",
					},
					To: v1.RouteTargetReference{
						Kind:   "Service",
						Name:   "svc1",
						Weight: &weight,
					},
				},
			},
			Want: `
        openshift_route_created{route="route1",namespace="ns1"} 1.5e+09
        openshift_route_labels{route="route1",label_app="good1",namespace="ns1"} 1
				openshift_route_info{route="route1",namespace="ns1",host="example.com",path="",tls_termination="edge",to_kind="Service",to_name="svc1",to_weight="100"} 1
				openshift_route_status{route="route1",namespace="ns1",host="example.com",status="True",type="Admitted",router_name="router1"} 1
				openshift_route_status{route="route1",namespace="ns1",host="example.com",status="True",type="Admitted",router_name="router2"} 1
				`,
			MetricNames: []string{"openshift_route_created", "openshift_route_labels", "openshift_route_info", "openshift_route_status"},
		},
	}

	for i, c := range cases {
		c.Func = metric.ComposeMetricGenFuncs(routeMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
