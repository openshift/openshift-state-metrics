package main

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/openshift/origin/pkg/util/proc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/openshift/openshift-state-metrics/pkg/version"
	kcollectors "k8s.io/kube-state-metrics/pkg/collector"
	koptions "k8s.io/kube-state-metrics/pkg/options"
	"k8s.io/kube-state-metrics/pkg/whiteblacklist"

	ocollectors "github.com/openshift/openshift-state-metrics/pkg/collectors"
	"github.com/openshift/openshift-state-metrics/pkg/options"
)

const (
	metricsPath = "/metrics"
	healthzPath = "/healthz"
)

// promLogger implements promhttp.Logger
type promLogger struct{}

func (pl promLogger) Println(v ...interface{}) {
	glog.Error(v...)
}

func main() {
	opts := options.NewOptions()
	opts.AddFlags()

	err := opts.Parse()
	if err != nil {
		glog.Fatalf("Error: %s", err)
	}

	if opts.Version {
		fmt.Printf("%#v\n", version.GetVersion())
		os.Exit(0)
	}

	if opts.Help {
		opts.Usage()
		os.Exit(0)
	}

	collectorBuilder := ocollectors.NewBuilder(context.TODO())
	collectorBuilder.WithApiserver(opts.Apiserver).WithKubeConfig(opts.Kubeconfig)
	if len(opts.Collectors) == 0 {
		glog.Info("Using default collectors")
		collectorBuilder.WithEnabledCollectors(options.DefaultCollectors.AsSlice())
	} else {
		collectorBuilder.WithEnabledCollectors(opts.Collectors.AsSlice())
	}

	if len(opts.Namespaces) == 0 {
		glog.Info("Using all namespace")
		collectorBuilder.WithNamespaces(koptions.DefaultNamespaces)
	} else {
		if opts.Namespaces.IsAllNamespaces() {
			glog.Info("Using all namespace")
		} else {
			glog.Infof("Using %s namespaces", opts.Namespaces)
		}
		collectorBuilder.WithNamespaces(opts.Namespaces)
	}

	whiteBlackList, err := whiteblacklist.New(opts.MetricWhitelist, opts.MetricBlacklist)
	if err != nil {
		glog.Fatal(err)
	}

	glog.Infof("metric white- blacklisting: %v", whiteBlackList.Status())

	collectorBuilder.WithWhiteBlackList(whiteBlackList)

	proc.StartReaper()

	if err != nil {
		glog.Fatalf("Failed to create client: %v", err)
	}

	osMetricsRegistry := prometheus.NewRegistry()
	osMetricsRegistry.Register(ocollectors.ResourcesPerScrapeMetric)
	osMetricsRegistry.Register(ocollectors.ScrapeErrorTotalMetric)
	osMetricsRegistry.Register(prometheus.NewProcessCollector(os.Getpid(), ""))
	osMetricsRegistry.Register(prometheus.NewGoCollector())
	go telemetryServer(osMetricsRegistry, opts.TelemetryHost, opts.TelemetryPort)

	collectors := collectorBuilder.Build()

	serveMetrics(collectors, opts.Host, opts.Port, opts.EnableGZIPEncoding)
}
func telemetryServer(registry prometheus.Gatherer, host string, port int) {
	// Address to listen on for web interface and telemetry
	listenAddress := net.JoinHostPort(host, strconv.Itoa(port))

	glog.Infof("Starting openshift-state-metrics self metrics server: %s", listenAddress)

	mux := http.NewServeMux()

	// Add metricsPath
	mux.Handle(metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{ErrorLog: promLogger{}}))
	// Add index
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>openshift-State-Metrics Metrics Server</title></head>
             <body>
             <h1>openshift-State-Metrics Metrics</h1>
			 <ul>
             <li><a href='` + metricsPath + `'>metrics</a></li>
			 </ul>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(listenAddress, mux))
}

// TODO: How about accepting an interface Collector instead?
func serveMetrics(collectors []*kcollectors.Collector, host string, port int, enableGZIPEncoding bool) {
	// Address to listen on for web interface and telemetry
	listenAddress := net.JoinHostPort(host, strconv.Itoa(port))

	glog.Infof("Starting metrics server: %s", listenAddress)

	mux := http.NewServeMux()

	// TODO: This doesn't belong into serveMetrics
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	// Add metricsPath
	mux.Handle(metricsPath, &metricHandler{collectors, enableGZIPEncoding})
	// Add healthzPath
	mux.HandleFunc(healthzPath, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})
	// Add index
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>OpenShift Metrics Server</title></head>
             <body>
             <h1>Kube Metrics</h1>
			 <ul>
             <li><a href='` + metricsPath + `'>metrics</a></li>
             <li><a href='` + healthzPath + `'>healthz</a></li>
			 </ul>
             </body>
             </html>`))
	})
	log.Fatal(http.ListenAndServe(listenAddress, mux))
}

type metricHandler struct {
	collectors         []*kcollectors.Collector
	enableGZIPEncoding bool
}

func (m *metricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resHeader := w.Header()
	var writer io.Writer = w

	resHeader.Set("Content-Type", `text/plain; version=`+"0.0.4")

	if m.enableGZIPEncoding {
		// Gzip response if requested. Taken from
		// github.com/prometheus/client_golang/prometheus/promhttp.decorateWriter.
		reqHeader := r.Header.Get("Accept-Encoding")
		parts := strings.Split(reqHeader, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "gzip" || strings.HasPrefix(part, "gzip;") {
				writer = gzip.NewWriter(writer)
				resHeader.Set("Content-Encoding", "gzip")
			}
		}
	}

	for _, c := range m.collectors {
		c.Collect(w)
	}

	// In case we gziped the response, we have to close the writer.
	if closer, ok := writer.(io.Closer); ok {
		closer.Close()
	}
}
