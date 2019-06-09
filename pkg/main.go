package main

import (
	"external-metrics/pkg/metrics_server"
	"flag"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/component-base/logs"
	"k8s.io/klog"
	"net/http"

	basecmd "github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/cmd"
)

type ExternalMetricsAdapter struct {
	basecmd.AdapterBase

	// Message is printed on succesful startup
	Message string
}


func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	cmd := &ExternalMetricsAdapter{}
	cmd.Flags().StringVar(&cmd.Message, "msg", "starting adapter...", "startup message")
	cmd.Flags().AddGoFlagSet(flag.CommandLine) // make sure we get the klog flags

	externalMetricsServer := metrics_server.NewServer()
	cmd.WithExternalMetrics(externalMetricsServer)

	klog.Infof(cmd.Message)

	go func() {
		// Open port for POSTing fake metrics
		klog.Fatal(http.ListenAndServe(":8080", nil))
	}()
	if err := cmd.Run(wait.NeverStop); err != nil {
		klog.Fatalf("unable to run external metrics adapter: %v", err)
	}
}
