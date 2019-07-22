package collectors

import (
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/openshift/api/quota/v1"
	"k8s.io/kube-state-metrics/pkg/metric"
)

func TestClusterResourceQuotaCollector(t *testing.T) {
	// Fixed metadata on type and help text. We prepend this to every expected
	// output so we only have to modify a single place when doing adjustments.
	const metadata = `
		# HELP openshift_clusterresourcequota_created Unix creation timestamp
		# TYPE openshift_clusterresourcequota_created gauge
		# HELP openshift_clusterresourcequota_labels Kubernetes labels converted to Prometheus labels.
		# TYPE openshift_clusterresourcequota_labels gauge
		# HELP openshift_clusterresourcequota_selector Selector of clusterresource quota, which defines the affected namespaces
		# TYPE openshift_clusterresourcequota_selector gauge
		# HELP openshift_clusterresourcequota_usage Usage about resource quota
		# TYPE openshift_clusterresourcequota_usage gauge
		# HELP openshift_clusterresourcequota_namespace_usage Usage about applied resource quota per namespace
		# TYPE openshift_clusterresourcequota_namespace_usage gauge
`
	cases := []generateMetricsTestCase{
		{
			Obj: &v1.ClusterResourceQuota{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "quota1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Labels: map[string]string{
						"quota": "test",
					},
				},
				Spec: v1.ClusterResourceQuotaSpec{},
				Status: v1.ClusterResourceQuotaStatus{
					Total: corev1.ResourceQuotaStatus{},
				},
			},
			Want: `
		openshift_clusterresourcequota_created{name="quota1"} 1.5e+09
        openshift_clusterresourcequota_labels{label_quota="test",name="quota1"} 1
`,
			MetricNames: []string{"openshift_clusterresourcequota_created", "openshift_clusterresourcequota_selector", "openshift_clusterresourcequota_labels", "openshift_clusterresourcequota_usage"},
		},

		{
			Obj: &v1.ClusterResourceQuota{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "quota1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Labels: map[string]string{
						"quota": "test",
					},
				},
				Spec: v1.ClusterResourceQuotaSpec{
					Selector: v1.ClusterResourceQuotaSelector{
						AnnotationSelector: map[string]string{},
						LabelSelector: &metav1.LabelSelector{
							// matchExpressions:
							//    - {key: tier, operator: In, values: [cache]}
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{
									Key:      "tier",
									Operator: "In",
									Values:   []string{"cache2,cache3"},
								},
								{
									Key:      "tier2",
									Operator: "DoesNotExist",
									//not set: Values:   []string{""},
								},
							},
						},
					},
				},
				Status: v1.ClusterResourceQuotaStatus{
					Total: corev1.ResourceQuotaStatus{},
				},
			},
			Want: `
		openshift_clusterresourcequota_created{name="quota1"} 1.5e+09
		openshift_clusterresourcequota_labels{label_quota="test",name="quota1"} 1
		openshift_clusterresourcequota_selector{name="quota1",type="match-expressions",key="tier",operator="In",values="cache2,cache3"} 1
		openshift_clusterresourcequota_selector{name="quota1",type="match-expressions",key="tier2",operator="DoesNotExist",values=""} 1
`,
			MetricNames: []string{"openshift_clusterresourcequota_created", "openshift_clusterresourcequota_selector", "openshift_clusterresourcequota_labels", "openshift_clusterresourcequota_usage"},
		},

		{
			Obj: &v1.ClusterResourceQuota{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "quota1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
				},
				Spec: v1.ClusterResourceQuotaSpec{
					Quota: corev1.ResourceQuotaSpec{
						Hard: corev1.ResourceList{
							corev1.ResourceCPU: resource.MustParse("4.3"),
						},
					},
					Selector: v1.ClusterResourceQuotaSelector{
						AnnotationSelector: map[string]string{},
						LabelSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{"clusterquota": "labeltest"},
						},
					},
				},
				Status: v1.ClusterResourceQuotaStatus{
					Total: corev1.ResourceQuotaStatus{},
				},
			},
			Want: `
		openshift_clusterresourcequota_created{name="quota1"} 1.5e+09
		openshift_clusterresourcequota_labels{name="quota1"} 1
        openshift_clusterresourcequota_selector{name="quota1",type="match-labels",key="clusterquota",value="labeltest"} 1
		openshift_clusterresourcequota_usage{name="quota1",resource="cpu",type="hard"} 4.3
`,

			MetricNames: []string{"openshift_clusterresourcequota_created", "openshift_clusterresourcequota_selector", "openshift_clusterresourcequota_labels", "openshift_clusterresourcequota_usage"},
		},
		{
			Obj: &v1.ClusterResourceQuota{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "quota1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Labels: map[string]string{
						"quota": "test",
					},
				},
				Spec: v1.ClusterResourceQuotaSpec{
					Quota: corev1.ResourceQuotaSpec{
						Hard: corev1.ResourceList{
							corev1.ResourceCPU:                    resource.MustParse("4.3"),
							corev1.ResourceMemory:                 resource.MustParse("2.1G"),
							corev1.ResourceStorage:                resource.MustParse("10G"),
							corev1.ResourcePods:                   resource.MustParse("9"),
							corev1.ResourceServices:               resource.MustParse("8"),
							corev1.ResourceReplicationControllers: resource.MustParse("7"),
							corev1.ResourceQuotas:                 resource.MustParse("6"),
							corev1.ResourceSecrets:                resource.MustParse("5"),
							corev1.ResourceConfigMaps:             resource.MustParse("4"),
							corev1.ResourcePersistentVolumeClaims: resource.MustParse("3"),
							corev1.ResourceServicesNodePorts:      resource.MustParse("2"),
							corev1.ResourceServicesLoadBalancers:  resource.MustParse("1"),
						},
					},
				},
				Status: v1.ClusterResourceQuotaStatus{
					Total: corev1.ResourceQuotaStatus{
						Hard: corev1.ResourceList{
							corev1.ResourceCPU:                    resource.MustParse("4.3"),
							corev1.ResourceMemory:                 resource.MustParse("2.1G"),
							corev1.ResourceStorage:                resource.MustParse("10G"),
							corev1.ResourcePods:                   resource.MustParse("9"),
							corev1.ResourceServices:               resource.MustParse("8"),
							corev1.ResourceReplicationControllers: resource.MustParse("7"),
							corev1.ResourceQuotas:                 resource.MustParse("6"),
							corev1.ResourceSecrets:                resource.MustParse("5"),
							corev1.ResourceConfigMaps:             resource.MustParse("4"),
							corev1.ResourcePersistentVolumeClaims: resource.MustParse("3"),
							corev1.ResourceServicesNodePorts:      resource.MustParse("2"),
							corev1.ResourceServicesLoadBalancers:  resource.MustParse("1"),
						},
						Used: corev1.ResourceList{
							corev1.ResourceCPU:                    resource.MustParse("2.1"),
							corev1.ResourceMemory:                 resource.MustParse("500M"),
							corev1.ResourceStorage:                resource.MustParse("9G"),
							corev1.ResourcePods:                   resource.MustParse("8"),
							corev1.ResourceServices:               resource.MustParse("7"),
							corev1.ResourceReplicationControllers: resource.MustParse("6"),
							corev1.ResourceQuotas:                 resource.MustParse("5"),
							corev1.ResourceSecrets:                resource.MustParse("4"),
							corev1.ResourceConfigMaps:             resource.MustParse("3"),
							corev1.ResourcePersistentVolumeClaims: resource.MustParse("2"),
							corev1.ResourceServicesNodePorts:      resource.MustParse("1"),
							corev1.ResourceServicesLoadBalancers:  resource.MustParse("0"),
						},
					},
				},
			},
			Want: `
       	openshift_clusterresourcequota_created{name="quota1"} 1.5e+09
        openshift_clusterresourcequota_labels{label_quota="test",name="quota1"} 1
        openshift_clusterresourcequota_usage{name="quota1",resource="configmaps",type="hard"} 4
        openshift_clusterresourcequota_usage{name="quota1",resource="configmaps",type="used"} 3
        openshift_clusterresourcequota_usage{name="quota1",resource="cpu",type="hard"} 4.3
        openshift_clusterresourcequota_usage{name="quota1",resource="cpu",type="used"} 2.1
        openshift_clusterresourcequota_usage{name="quota1",resource="memory",type="hard"} 2.1e+09
        openshift_clusterresourcequota_usage{name="quota1",resource="memory",type="used"} 5e+08
        openshift_clusterresourcequota_usage{name="quota1",resource="persistentvolumeclaims",type="hard"} 3
        openshift_clusterresourcequota_usage{name="quota1",resource="persistentvolumeclaims",type="used"} 2
        openshift_clusterresourcequota_usage{name="quota1",resource="pods",type="hard"} 9
        openshift_clusterresourcequota_usage{name="quota1",resource="pods",type="used"} 8
        openshift_clusterresourcequota_usage{name="quota1",resource="replicationcontrollers",type="hard"} 7
        openshift_clusterresourcequota_usage{name="quota1",resource="replicationcontrollers",type="used"} 6
        openshift_clusterresourcequota_usage{name="quota1",resource="resourcequotas",type="hard"} 6
        openshift_clusterresourcequota_usage{name="quota1",resource="resourcequotas",type="used"} 5
        openshift_clusterresourcequota_usage{name="quota1",resource="secrets",type="hard"} 5
        openshift_clusterresourcequota_usage{name="quota1",resource="secrets",type="used"} 4
        openshift_clusterresourcequota_usage{name="quota1",resource="services",type="hard"} 8
        openshift_clusterresourcequota_usage{name="quota1",resource="services",type="used"} 7
        openshift_clusterresourcequota_usage{name="quota1",resource="services.loadbalancers",type="hard"} 1
        openshift_clusterresourcequota_usage{name="quota1",resource="services.loadbalancers",type="used"} 0
        openshift_clusterresourcequota_usage{name="quota1",resource="services.nodeports",type="hard"} 2
        openshift_clusterresourcequota_usage{name="quota1",resource="services.nodeports",type="used"} 1
        openshift_clusterresourcequota_usage{name="quota1",resource="storage",type="hard"} 1e+10
        openshift_clusterresourcequota_usage{name="quota1",resource="storage",type="used"} 9e+09
`,

			MetricNames: []string{"openshift_clusterresourcequota_created", "openshift_clusterresourcequota_labels", "openshift_clusterresourcequota"},
		},
		{
			Obj: &v1.ClusterResourceQuota{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "quota1",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
					Namespace:         "ns1",
					Labels: map[string]string{
						"quota": "test",
					},
				},
				Spec: v1.ClusterResourceQuotaSpec{
					Quota: corev1.ResourceQuotaSpec{
						Hard: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse("2.1G"),
						},
					},
					Selector: v1.ClusterResourceQuotaSelector{
						AnnotationSelector: map[string]string{
							"clusterquota": "test",
						},
						LabelSelector: &metav1.LabelSelector{},
					},
				},
				Status: v1.ClusterResourceQuotaStatus{
					Namespaces: []v1.ResourceQuotaStatusByNamespace{
						{
							Namespace: "myproject",
							Status: corev1.ResourceQuotaStatus{
								Hard: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("2.1G"),
								},
								Used: corev1.ResourceList{
									corev1.ResourceMemory: resource.MustParse("500M"),
								},
							},
						},
					},
					Total: corev1.ResourceQuotaStatus{
						Hard: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse("2.1G"),
						},
						Used: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse("500M"),
						},
					},
				},
			},
			Want: `
       	openshift_clusterresourcequota_created{name="quota1"} 1.5e+09
        openshift_clusterresourcequota_selector{name="quota1",type="annotation",key="clusterquota",value="test"} 1
	    openshift_clusterresourcequota_usage{name="quota1",resource="memory",type="hard"} 2.1e+09
        openshift_clusterresourcequota_usage{name="quota1",resource="memory",type="used"} 5e+08		
		openshift_clusterresourcequota_namespace_usage{name="quota1",namespace="myproject",resource="memory",type="hard"} 2.1e+09
        openshift_clusterresourcequota_namespace_usage{name="quota1",namespace="myproject",resource="memory",type="used"} 5e+08		
`,

			MetricNames: []string{"openshift_clusterresourcequota_created", "openshift_clusterresourcequota_selector", "openshift_clusterresourcequota_usage", "openshift_clusterresourcequota_namespace_usage"},
		},
	}

	for i, c := range cases {
		c.Func = metric.ComposeMetricGenFuncs(quotaMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
