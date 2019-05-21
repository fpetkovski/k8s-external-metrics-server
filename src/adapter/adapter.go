/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package adapter

import (
	"external-metrics/beanstalkd_client"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/dynamic"
	"k8s.io/metrics/pkg/apis/external_metrics"
	"sync"
	"time"

	"github.com/kubernetes-incubator/custom-metrics-apiserver/pkg/provider"
)

type externalMetric struct {
	info   provider.ExternalMetricInfo
	labels map[string]string
	value  external_metrics.ExternalMetricValue
}

var (
	defaultExternalMetrics = []externalMetric{
		makeMetric("default", 0),
	}
)

func makeMetric(metricName string, metricValue int64) externalMetric {
	return externalMetric{
		info: provider.ExternalMetricInfo{
			Metric: "tube-" + metricName,
		},
		labels: map[string]string{"tube": metricName},
		value: external_metrics.ExternalMetricValue{
			MetricName: "tube-" + metricName,
			MetricLabels: map[string]string{
				"tube": metricName,
			},
			Value: *resource.NewQuantity(metricValue, resource.DecimalSI),
		},
	}
}

type metricsProvider struct {
	client dynamic.Interface
	mapper apimeta.RESTMapper

	valuesLock      sync.RWMutex
	externalMetrics []externalMetric
}

func pollMetrics(provider *metricsProvider)  {
	beanstalkdClient := beanstalkd_client.NewClient("beanstalkd", "11300")

	for {
		totalJobs := beanstalkdClient.GetJobsCount("default")
		provider.externalMetrics = []externalMetric{
			makeMetric("default", totalJobs),
		}
		time.Sleep(5 * time.Second)
	}
}

func NewProvider(client dynamic.Interface, mapper apimeta.RESTMapper) provider.ExternalMetricsProvider {
	provider := &metricsProvider{
		client:          client,
		mapper:          mapper,
		externalMetrics: defaultExternalMetrics,
	}

	go pollMetrics(provider)

	return provider
}

func (p *metricsProvider) GetExternalMetric(namespace string, metricSelector labels.Selector, info provider.ExternalMetricInfo) (*external_metrics.ExternalMetricValueList, error) {
	p.valuesLock.RLock()
	defer p.valuesLock.RUnlock()

	matchingMetrics := []external_metrics.ExternalMetricValue{}
	for _, metric := range p.externalMetrics {
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

func (p *metricsProvider) ListAllExternalMetrics() []provider.ExternalMetricInfo {
	p.valuesLock.RLock()
	defer p.valuesLock.RUnlock()

	externalMetricsInfo := []provider.ExternalMetricInfo{}
	for _, metric := range p.externalMetrics {
		externalMetricsInfo = append(externalMetricsInfo, metric.info)
	}
	return externalMetricsInfo
}
