package options

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	koptions "k8s.io/kube-state-metrics/pkg/options"
)

func init() {
	//TODO this is because the CollectorSet struct is validate the collectors from the commandline using
	//"DefaultCollectors". https://github.com/kubernetes/kube-state-metrics/blob/master/pkg/options/types.go#L80
	koptions.DefaultCollectors["deploymentConfigs"] = struct{}{}
	koptions.DefaultCollectors["buildconfigs"] = struct{}{}
	koptions.DefaultCollectors["builds"] = struct{}{}
	koptions.DefaultCollectors["clusterresourcequotas"] = struct{}{}
}

var (
	DefaultNamespaces = koptions.NamespaceList{metav1.NamespaceAll}
	DefaultCollectors = koptions.CollectorSet{
		"deploymentConfigs":     struct{}{},
		"buildconfigs":          struct{}{},
		"builds":                struct{}{},
		"clusterresourcequotas": struct{}{},
	}
)
