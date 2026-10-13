package collectors

import (
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/openshift/api/apps/v1"

	"k8s.io/kube-state-metrics/pkg/metric"
)

var (
	depl1Replicas int32 = 200
	depl2Replicas int32 = 5

	depl1MaxUnavailable = intstr.FromInt(10)
	depl2MaxUnavailable = intstr.FromString("20%")

	depl1MaxSurge = intstr.FromInt(10)
	depl2MaxSurge = intstr.FromString("20%")
)

func TestDeploymentCollector(t *testing.T) {
	// Fixed metadata on type and help text. We prepend this to every expected
	// output so we only have to modify a single place when doing adjustments.
	const metadata = `
		# HELP openshift_deploymentconfig_created Unix creation timestamp
		# TYPE openshift_deploymentconfig_created gauge
		# HELP openshift_deploymentconfig_metadata_generation Sequence number representing a specific generation of the desired state.
		# TYPE openshift_deploymentconfig_metadata_generation gauge
		# HELP openshift_deploymentconfig_spec_paused Whether the deployment is paused and will not be processed by the deployment controller.
		# TYPE openshift_deploymentconfig_spec_paused gauge
		# HELP openshift_deploymentconfig_spec_replicas Number of desired pods for a deployment.
		# TYPE openshift_deploymentconfig_spec_replicas gauge
		# HELP openshift_deploymentconfig_status_replicas The number of replicas per deployment.
		# TYPE openshift_deploymentconfig_status_replicas gauge
		# HELP openshift_deploymentconfig_status_replicas_available The number of available replicas per deployment.
		# TYPE openshift_deploymentconfig_status_replicas_available gauge
		# HELP openshift_deploymentconfig_status_replicas_unavailable The number of unavailable replicas per deployment.
		# TYPE openshift_deploymentconfig_status_replicas_unavailable gauge
		# HELP openshift_deploymentconfig_status_replicas_updated The number of updated replicas per deployment.
		# TYPE openshift_deploymentconfig_status_replicas_updated gauge
		# HELP openshift_deploymentconfig_status_observed_generation The generation observed by the deployment controller.
		# TYPE openshift_deploymentconfig_status_observed_generation gauge
		# HELP openshift_deploymentconfig_spec_strategy_rollingupdate_max_unavailable Maximum number of unavailable replicas during a rolling update of a deployment.
		# TYPE openshift_deploymentconfig_spec_strategy_rollingupdate_max_unavailable gauge
		# HELP openshift_deploymentconfig_spec_strategy_rollingupdate_max_surge Maximum number of replicas that can be scheduled above the desired number of replicas during a rolling update of a deployment.
		# TYPE openshift_deploymentconfig_spec_strategy_rollingupdate_max_surge gauge
		# HELP openshift_deploymentconfig_labels Kubernetes labels converted to Prometheus labels.
		# TYPE openshift_deploymentconfig_labels gauge
	`
	cases := []generateMetricsTestCase{
		{
			Obj: &v1.DeploymentConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "depl1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Labels: map[string]string{
						"app": "example1",
					},
					Generation: 21,
				},
				Status: v1.DeploymentConfigStatus{
					Replicas:            15,
					AvailableReplicas:   10,
					UnavailableReplicas: 5,
					UpdatedReplicas:     2,
					ObservedGeneration:  111,
				},
				Spec: v1.DeploymentConfigSpec{
					Replicas: depl1Replicas,
					Strategy: v1.DeploymentStrategy{
						RollingParams: &v1.RollingDeploymentStrategyParams{
							MaxUnavailable: &depl1MaxUnavailable,
							MaxSurge:       &depl1MaxSurge,
						},
					},
				},
			},
			Want: `
        openshift_deploymentconfig_created{deploymentconfig="depl1",namespace="ns1"} 1.5e+09
        openshift_deploymentconfig_labels{deploymentconfig="depl1",label_app="example1",namespace="ns1"} 1
        openshift_deploymentconfig_metadata_generation{deploymentconfig="depl1",namespace="ns1"} 21
        openshift_deploymentconfig_spec_paused{deploymentconfig="depl1",namespace="ns1"} 0
        openshift_deploymentconfig_spec_replicas{deploymentconfig="depl1",namespace="ns1"} 200
        openshift_deploymentconfig_spec_strategy_rollingupdate_max_surge{deploymentconfig="depl1",namespace="ns1"} 10
        openshift_deploymentconfig_spec_strategy_rollingupdate_max_unavailable{deploymentconfig="depl1",namespace="ns1"} 10
        openshift_deploymentconfig_status_observed_generation{deploymentconfig="depl1",namespace="ns1"} 111
        openshift_deploymentconfig_status_replicas_available{deploymentconfig="depl1",namespace="ns1"} 10
        openshift_deploymentconfig_status_replicas_unavailable{deploymentconfig="depl1",namespace="ns1"} 5
        openshift_deploymentconfig_status_replicas_updated{deploymentconfig="depl1",namespace="ns1"} 2
        openshift_deploymentconfig_status_replicas{deploymentconfig="depl1",namespace="ns1"} 15
`,
		},
		{
			Obj: &v1.DeploymentConfig{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "depl2",
					Namespace: "ns2",
					Labels: map[string]string{
						"app": "example2",
					},
					Generation: 14,
				},
				Status: v1.DeploymentConfigStatus{
					Replicas:            10,
					AvailableReplicas:   5,
					UnavailableReplicas: 0,
					UpdatedReplicas:     1,
					ObservedGeneration:  1111,
				},
				Spec: v1.DeploymentConfigSpec{
					Paused:   true,
					Replicas: depl2Replicas,
					Strategy: v1.DeploymentStrategy{
						RollingParams: &v1.RollingDeploymentStrategyParams{
							MaxUnavailable: &depl2MaxUnavailable,
							MaxSurge:       &depl2MaxSurge,
						},
					},
				},
			},
			Want: `
       	openshift_deploymentconfig_labels{deploymentconfig="depl2",label_app="example2",namespace="ns2"} 1
        openshift_deploymentconfig_metadata_generation{deploymentconfig="depl2",namespace="ns2"} 14
        openshift_deploymentconfig_spec_paused{deploymentconfig="depl2",namespace="ns2"} 1
        openshift_deploymentconfig_spec_replicas{deploymentconfig="depl2",namespace="ns2"} 5
        openshift_deploymentconfig_spec_strategy_rollingupdate_max_surge{deploymentconfig="depl2",namespace="ns2"} 1
        openshift_deploymentconfig_spec_strategy_rollingupdate_max_unavailable{deploymentconfig="depl2",namespace="ns2"} 1
        openshift_deploymentconfig_status_observed_generation{deploymentconfig="depl2",namespace="ns2"} 1111
        openshift_deploymentconfig_status_replicas_available{deploymentconfig="depl2",namespace="ns2"} 5
        openshift_deploymentconfig_status_replicas_unavailable{deploymentconfig="depl2",namespace="ns2"} 0
        openshift_deploymentconfig_status_replicas_updated{deploymentconfig="depl2",namespace="ns2"} 1
        openshift_deploymentconfig_status_replicas{deploymentconfig="depl2",namespace="ns2"} 10
`,
		},
	}

	for i, c := range cases {
		c.Func = metric.ComposeMetricGenFuncs(deploymentMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
