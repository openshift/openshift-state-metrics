module github.com/openshift/openshift-state-metrics

go 1.14

require (
	github.com/openshift/api v0.0.0-20200623075207-eb651a5bb0ad
	github.com/openshift/client-go v0.0.0-20200623090625-83993cebb5ae
	github.com/prometheus/client_golang v1.7.1
	github.com/spf13/pflag v1.0.5
	golang.org/x/net v0.0.0-20210224082022-3d97a244fca7
	k8s.io/api v0.22.0-alpha.2
	k8s.io/apimachinery v0.22.0-alpha.2
	k8s.io/client-go v0.22.0-alpha.2
	k8s.io/klog/v2 v2.8.0
	k8s.io/kube-state-metrics v0.0.0-20190129120824-7bfed92869b6
)
