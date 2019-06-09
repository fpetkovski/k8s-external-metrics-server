package adapter

import (
	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/metrics/pkg/apis/external_metrics"
)

type externalMetric struct {
	info   provider.ExternalMetricInfo
	labels map[string]string
	value  external_metrics.ExternalMetricValue
}

func makeMetric(metricName string, metricValue int64) externalMetric {
	tubeName := "tube-" + metricName

	return externalMetric{
		info: provider.ExternalMetricInfo{
			Metric: tubeName,
		},
		labels: map[string]string{"tube": metricName},
		value: external_metrics.ExternalMetricValue{
			MetricName: tubeName,
			MetricLabels: map[string]string{
				"tube": metricName,
			},
			Value: *resource.NewQuantity(metricValue, resource.DecimalSI),
		},
	}
}
