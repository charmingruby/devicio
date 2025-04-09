package metric

import (
	"github.com/charmingruby/devicio/lib/observability"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusMeter struct {
	metrics map[string]prometheus.Collector
}

func NewPrometheusMeter() *PrometheusMeter {
	return &PrometheusMeter{
		metrics: make(map[string]prometheus.Collector),
	}
}

func (p *PrometheusMeter) NewHistogram(input observability.HistogramInput) {
	histogram := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      input.Name,
		Help:      input.Help,
		Namespace: input.Namespace,
		Buckets:   input.Buckets,
	})

	p.metrics[input.Name] = histogram

	prometheus.MustRegister(histogram)
}

func (p *PrometheusMeter) NewCounter(input observability.CounterInput) {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name:      input.Name,
		Help:      input.Help,
		Namespace: input.Namespace,
	})

	p.metrics[input.Name] = counter

	prometheus.MustRegister(counter)
}

func (p *PrometheusMeter) NewCounterList(input observability.CounterListInput) {
	counterList := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      input.Name,
		Help:      input.Help,
		Namespace: input.Namespace,
	}, input.LabelNames)

	p.metrics[input.Name] = counterList

	prometheus.MustRegister(counterList)
}

func (p *PrometheusMeter) GetMetric(name, metricType string) (any, error) {
	if !observability.ValidateMetricType(metricType) {
		return nil, observability.ErrInvalidMetricType
	}

	metric := p.metrics[name]
	if metric == nil {
		return nil, observability.ErrMetricNotFound
	}

	return metric, nil
}
