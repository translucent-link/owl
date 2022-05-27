package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ReqProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "tlscoring_processed_reqs_total",
		Help: "The total number of processed requests",
	})
)
