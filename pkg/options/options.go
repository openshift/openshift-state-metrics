package options

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/pflag"

	koptions "k8s.io/kube-state-metrics/pkg/options"
)

type Options struct {
	Apiserver       string
	Kubeconfig      string
	Help            bool
	Port            int
	Host            string
	TelemetryPort   int
	TelemetryHost   string
	Collectors      koptions.CollectorSet
	Namespaces      koptions.NamespaceList
	MetricBlacklist koptions.MetricSet
	MetricWhitelist koptions.MetricSet
	Version         bool

	EnableGZIPEncoding bool

	flags *pflag.FlagSet
}

func NewOptions() *Options {
	return &Options{
		Collectors:      koptions.CollectorSet{},
		MetricWhitelist: koptions.MetricSet{},
		MetricBlacklist: koptions.MetricSet{},
	}
}

func (o *Options) AddFlags() {
	o.flags = pflag.NewFlagSet("", pflag.ExitOnError)
	// add glog flags
	o.flags.AddGoFlagSet(flag.CommandLine)
	o.flags.Lookup("logtostderr").Value.Set("true")
	o.flags.Lookup("logtostderr").DefValue = "true"
	o.flags.Lookup("logtostderr").NoOptDefVal = "true"

	o.flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		o.flags.PrintDefaults()
	}

	o.flags.StringVar(&o.Apiserver, "apiserver", "", `The URL of the apiserver to use as a master`)
	o.flags.StringVar(&o.Kubeconfig, "kubeconfig", "", "Absolute path to the kubeconfig file")
	o.flags.BoolVarP(&o.Help, "help", "h", false, "Print Help text")
	o.flags.IntVar(&o.Port, "port", 80, `Port to expose metrics on.`)
	o.flags.StringVar(&o.Host, "host", "0.0.0.0", `Host to expose metrics on.`)
	o.flags.IntVar(&o.TelemetryPort, "telemetry-port", 81, `Port to expose openshift-state-metrics self metrics on.`)
	o.flags.StringVar(&o.TelemetryHost, "telemetry-host", "0.0.0.0", `Host to expose openshift-state-metrics self metrics on.`)
	o.flags.Var(&o.Collectors, "collectors", fmt.Sprintf("Comma-separated list of collectors to be enabled. Defaults to %q", &DefaultCollectors))
	o.flags.Var(&o.Namespaces, "namespace", fmt.Sprintf("Comma-separated list of namespaces to be enabled. Defaults to %q", &DefaultNamespaces))
	o.flags.Var(&o.MetricWhitelist, "metric-whitelist", "Comma-separated list of metrics to be exposed. The whitelist and blacklist are mutually exclusive.")
	o.flags.Var(&o.MetricBlacklist, "metric-blacklist", "Comma-separated list of metrics not to be enabled. The whitelist and blacklist are mutually exclusive.")
	o.flags.BoolVarP(&o.Version, "version", "", false, "openshift-state-metrics build version information")

	o.flags.BoolVar(&o.EnableGZIPEncoding, "enable-gzip-encoding", false, "Gzip responses when requested by clients via 'Accept-Encoding: gzip' header.")
}

func (o *Options) Parse() error {
	err := o.flags.Parse(os.Args)
	return err
}

func (o *Options) Usage() {
	o.flags.Usage()
}
