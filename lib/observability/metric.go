package observability

import "errors"

const (
	HistogramMetricType   = "histogram"
	CounterMetricType     = "counter"
	CounterListMetricType = "counter_list"
)

var (
	ErrInvalidMetricType = errors.New("invalid metric type")
	ErrMetricNotFound    = errors.New("metric not found")
)

type HistogramInput struct {
	Name      string
	Help      string
	Namespace string
	Buckets   []float64
}

type CounterInput struct {
	Name      string
	Help      string
	Namespace string
}

type CounterListInput struct {
	CounterInput
	LabelNames []string
}

type Meter interface {
	NewHistogram(input HistogramInput)
	NewCounter(input CounterInput)
	NewCounterList(input CounterListInput)
	GetMetric(name, metricType string) (any, error)
}

func ValidateMetricType(metricType string) bool {
	validTypes := map[string]bool{
		HistogramMetricType:   true,
		CounterMetricType:     true,
		CounterListMetricType: true,
	}

	return validTypes[metricType]
}
