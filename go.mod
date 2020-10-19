module github.com/openshift/openshift-state-metrics

go 1.14

require (
	github.com/openshift/api v0.0.0-20200623075207-eb651a5bb0ad
	github.com/openshift/client-go v0.0.0-20200623090625-83993cebb5ae
	github.com/openshift/origin v4.1.0+incompatible
	github.com/prometheus/client_golang v1.7.1
	github.com/spf13/pflag v1.0.5
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/klog/v2 v2.3.0
	k8s.io/kube-state-metrics v0.0.0-20190129120824-7bfed92869b6
)
