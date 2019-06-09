package adapter

import (
	"external-metrics/pkg/beanstalkd_client"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/dynamic"
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

type metricsProvider struct {
	client dynamic.Interface
	mapper apimeta.RESTMapper

	valuesLock      sync.RWMutex
	externalMetrics []externalMetric
}

func NewProvider(client dynamic.Interface, mapper apimeta.RESTMapper) provider.ExternalMetricsProvider {
	provider := &metricsProvider{
		client:          client,
		mapper:          mapper,
		externalMetrics: initialMetrics,
	}

	go pollMetrics(provider)

	return provider
}

func (provider *metricsProvider) GetExternalMetric(namespace string, metricSelector labels.Selector, info provider.ExternalMetricInfo) (*external_metrics.ExternalMetricValueList, error) {
	provider.valuesLock.RLock()
	defer provider.valuesLock.RUnlock()

	var matchingMetrics []external_metrics.ExternalMetricValue
	for _, metric := range provider.externalMetrics {
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

func (provider *metricsProvider) ListAllExternalMetrics() []provider.ExternalMetricInfo {
	provider.valuesLock.RLock()
	defer provider.valuesLock.RUnlock()

	var externalMetricsInfo []provider.ExternalMetricInfo
	for _, metric := range provider.externalMetrics {
		externalMetricsInfo = append(externalMetricsInfo, metric.info)
	}
	return externalMetricsInfo
}


func pollMetrics(provider *metricsProvider)  {
	beanstalkdClient := beanstalkd_client.NewClient("beanstalkd", "11300")

	for {
		totalJobs := beanstalkdClient.GetJobsCount("default")
		provider.updateTotalJobs(totalJobs)
		time.Sleep(5 * time.Second)
	}
}

func (provider *metricsProvider) updateTotalJobs(totalJobs int64) {
	provider.valuesLock.RLock()
	defer provider.valuesLock.RUnlock()

	provider.externalMetrics = []externalMetric{
		makeMetric("default", totalJobs),
	}
}