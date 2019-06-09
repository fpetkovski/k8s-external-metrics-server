package metrics_server

import (
	"external-metrics/pkg/beanstalkd_client"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/metrics/pkg/apis/external_metrics"
	"sync"
	"time"

	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
)

var (
	initialMetrics = []externalMetric{
		makeMetric("default", 0),
	}
)

type server struct {
	mapper apimeta.RESTMapper

	valuesLock      sync.RWMutex
	externalMetrics []externalMetric
}

func NewServer(mapper apimeta.RESTMapper) provider.ExternalMetricsProvider {
	server := &server{
		mapper:          mapper,
		externalMetrics: initialMetrics,
	}

	go pollMetrics(server)

	return server
}

func (server *server) GetExternalMetric(namespace string, metricSelector labels.Selector, info provider.ExternalMetricInfo) (*external_metrics.ExternalMetricValueList, error) {
	server.valuesLock.RLock()
	defer server.valuesLock.RUnlock()

	var matchingMetrics []external_metrics.ExternalMetricValue
	for _, metric := range server.externalMetrics {
		if metric.info.Metric == info.Metric &&
			metricSelector.Matches(labels.Set(metric.labels)) {
			metricValue := metric.value
			metricValue.Timestamp = metav1.Now()
			matchingMetrics = append(matchingMetrics, metricValue)
		}
	}
	return &external_metrics.ExternalMetricValueList{
		Items: matchingMetrics,
	}, nil
}

func (server *server) ListAllExternalMetrics() []provider.ExternalMetricInfo {
	server.valuesLock.RLock()
	defer server.valuesLock.RUnlock()

	var externalMetricsInfo []provider.ExternalMetricInfo
	for _, metric := range server.externalMetrics {
		externalMetricsInfo = append(externalMetricsInfo, metric.info)
	}
	return externalMetricsInfo
}

func pollMetrics(server *server) {
	beanstalkdClient := beanstalkd_client.NewClient("beanstalkd", "11300")

	for {
		totalJobs := beanstalkdClient.GetJobsCount("default")
		server.updateTotalJobs(totalJobs)
		time.Sleep(5 * time.Second)
	}
}

func (server *server) updateTotalJobs(totalJobs int64) {
	server.valuesLock.RLock()
	defer server.valuesLock.RUnlock()

	server.externalMetrics = []externalMetric{
		makeMetric("default", totalJobs),
	}
}
