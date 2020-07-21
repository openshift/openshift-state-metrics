module github.com/openshift/openshift-state-metrics

go 1.14

require (
	github.com/Azure/go-autorest/autorest v0.11.2 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/gophercloud/gophercloud v0.12.0 // indirect
	github.com/imdario/mergo v0.3.10 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/openshift/api v0.0.0-20200623075207-eb651a5bb0ad
	github.com/openshift/client-go v0.0.0-20200623090625-83993cebb5ae
	github.com/openshift/origin v4.1.0+incompatible
	github.com/prometheus/client_golang v0.0.0-20170531130054-e7e903064f5e
	github.com/prometheus/client_model v0.0.0-20150212101744-fa8ad6fec335 // indirect
	github.com/prometheus/common v0.0.0-20170427095455-13ba4ddd0caa // indirect
	github.com/prometheus/procfs v0.0.0-20170519190837-65c1f6f8f0fc // indirect
	github.com/spf13/pflag v1.0.5
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/kube-state-metrics v0.0.0-20190129120824-7bfed92869b6
	k8s.io/utils v0.0.0-20200720150651-0bdb4ca86cbc // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.18.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.3
	k8s.io/client-go => k8s.io/client-go v0.18.3
)
