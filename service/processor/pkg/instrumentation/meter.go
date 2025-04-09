package instrumentation

import (
	"net/http"

	"github.com/charmingruby/devicio/lib/observability"
	"github.com/charmingruby/devicio/lib/observability/metric"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Meter observability.Meter
)

func NewMeter() {
	Meter = metric.NewPrometheusMeter()

	Meter.NewCounter(observability.CounterInput{
		Name:      "messages_processed",
		Help:      "Total number of messages processed",
		Namespace: "devicio",
	})

	Meter.NewHistogram(observability.HistogramInput{
		Name:      "processing_time",
		Help:      "Time taken to process messages in seconds",
		Namespace: "devicio",
	})

	Meter.NewCounterList(observability.CounterListInput{
		CounterInput: observability.CounterInput{
			Name:      "errors",
			Help:      "Total number of errors by type",
			Namespace: "devicio",
		},
		LabelNames: []string{"error_type"},
	})
}

func RunMetricsServer(port string) error {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		return err
	}

	return nil
}

func ProcessingTimeHistogramMetric() (*prometheus.HistogramVec, error) {
	metric, err := Meter.GetMetric("processing_time", observability.HistogramMetricType)
	if err != nil {
		return nil, err
	}

	return metric.(*prometheus.HistogramVec), nil
}

func MessagesProcessedCounterMetric() (prometheus.Counter, error) {
	metric, err := Meter.GetMetric("messages_processed", observability.CounterMetricType)
	if err != nil {
		return nil, err
	}

	return metric.(prometheus.Counter), nil
}

func ErrorsCounterListMetric() (*prometheus.CounterVec, error) {
	metric, err := Meter.GetMetric("errors", observability.CounterListMetricType)
	if err != nil {
		return nil, err
	}

	return metric.(*prometheus.CounterVec), nil
}
