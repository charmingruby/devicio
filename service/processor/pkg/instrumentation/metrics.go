package instrumentation

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MessagesProcessedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "devicio_messages_processed_total",
		Help: "Total number of messages processed",
	})

	ProcessingTimeHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "devicio_processing_time_seconds",
		Help:    "Time taken to process messages in seconds",
		Buckets: prometheus.ExponentialBuckets(0.1, 2, 5),
	})

	ErrorCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "devicio_errors_total",
		Help: "Total number of errors by type",
	}, []string{"error_type"})

	QueueLatencyGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "devicio_queue_latency_seconds",
		Help: "Current queue latency in seconds",
	})
)
