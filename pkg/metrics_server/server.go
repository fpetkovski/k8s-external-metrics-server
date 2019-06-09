package metrics_server

import (
	"external-metrics/pkg/beanstalkd_client"
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

type queueInterface interface {
	GetJobsCount(tubeName string) int64
}

type server struct {
	valuesLock      sync.RWMutex
	queue           queueInterface
	externalMetrics []externalMetric
}

func NewServer() provider.ExternalMetricsProvider {
	beanstalkdClient := beanstalkd_client.NewClient("beanstalkd", "11300")

	server := &server{
		externalMetrics: initialMetrics,
		queue:           beanstalkdClient,
	}

	go server.pollMetrics()

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

func (server *server) pollMetrics() {
	for {
		totalJobs := server.queue.GetJobsCount("default")
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
