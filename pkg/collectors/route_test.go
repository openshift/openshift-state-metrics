package collectors

import (
	"testing"

	v1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kube-state-metrics/pkg/metric"
)

var (
	weight100 int32 = 100
	weight0   int32 = 0
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
		# HELP openshift_route_status Information about route status.
		# TYPE openshift_route_status gauge
	`
	cases := []generateMetricsTestCase{

		{
			Obj: &v1.Route{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route1",
					Namespace: "ns1",
				},
				Spec: v1.RouteSpec{
					Host: "example.com",
					TLS: &v1.TLSConfig{
						Termination: "edge",
					},
					To: v1.RouteTargetReference{
						Kind:   "Service",
						Name:   "svc1",
						Weight: &weight100,
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
								{
									Type:   v1.RouteAdmitted,
									Status: corev1.ConditionTrue,
								},
							},
						},
					},
				},
			},
			Want: `
		openshift_route_info{host="example.com",namespace="ns1",path="",route="route1",tls_termination="edge",to_kind="Service",to_name="svc1",to_weight="100"} 1
		openshift_route_status{route="route1",namespace="ns1",host="example.com",status="True",type="Admitted",router_name="router1"} 1
		openshift_route_status{route="route1",namespace="ns1",host="example.com",status="True",type="Admitted",router_name="router2"} 1
				`,
			MetricNames: []string{"openshift_route_info", "openshift_route_created", "openshift_route_status"},
		},
		{
			Obj: &v1.Route{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "route1",
					Namespace: "ns1",
				},
				Spec: v1.RouteSpec{
					Host: "example.com",
					TLS: &v1.TLSConfig{
						Termination: "edge",
					},
					To: v1.RouteTargetReference{
						Kind:   "Service",
						Name:   "svc1",
						Weight: &weight0,
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
								{
									Type:   v1.RouteAdmitted,
									Status: corev1.ConditionFalse,
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
								{
									Type:   v1.RouteAdmitted,
									Status: corev1.ConditionTrue,
								},
							},
						},
					},
				},
			},
			Want: `
		openshift_route_info{host="example.com",namespace="ns1",path="",route="route1",tls_termination="edge",to_kind="Service",to_name="svc1",to_weight="0"} 1
		openshift_route_status{route="route1",namespace="ns1",host="example.com",status="False",type="Admitted",router_name="router1"} 1
		openshift_route_status{route="route1",namespace="ns1",host="example.com",status="True",type="Admitted",router_name="router1"} 1
		openshift_route_status{route="route1",namespace="ns1",host="example.com",status="True",type="Admitted",router_name="router2"} 1
				`,
			MetricNames: []string{"openshift_route_info", "openshift_route_created", "openshift_route_status"},
		},
	}

	for i, c := range cases {
		c.Func = metric.ComposeMetricGenFuncs(routeMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
