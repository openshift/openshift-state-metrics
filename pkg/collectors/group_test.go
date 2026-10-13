package collectors

import (
	"testing"
	"time"

	"github.com/openshift/api/user/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kube-state-metrics/pkg/metric"
)

func TestGroupCollector(t *testing.T) {
	// Fixed metadata on type and help text. We prepend this to every expected
	// output so we only have to modify a single place when doing adjustments.
	const metadata = `
		# HELP openshift_group_created Unix creation timestamp
		# TYPE openshift_group_created gauge
		# HELP openshift_group_user_account User account in a group.
		# TYPE openshift_group_user_account gauge
	`
	cases := []generateMetricsTestCase{
		{
			Obj: &v1.Group{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "group",
					CreationTimestamp: metav1.Time{Time: time.Unix(1500000000, 0)},
				},
				Users: []string{"user1"},
			},
			Want: `
        openshift_group_created{group="group"} 1.5e+09
        openshift_group_user_account{group="group",user="user1"} 1
				`,
			MetricNames: []string{"openshift_group_created", "openshift_group_user_account"},
		},
	}

	for i, c := range cases {
		c.Func = metric.ComposeMetricGenFuncs(groupMetricFamilies)
		if err := c.run(); err != nil {
			t.Errorf("unexpected collecting result in %vth run:\n%s", i, err)
		}
	}
}
