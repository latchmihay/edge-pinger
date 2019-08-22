package prom

import "github.com/prometheus/client_golang/prometheus"

func AddGauge(subsystem string, gaugeName string, description string, labels []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: subsystem,
			Name:      gaugeName,
			Help:      description,
		},
		labels,
	)
}

func AddCounter(subsystem string, counterName string, description string, labels []string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      counterName,
			Help:      description,
		},
		labels,
	)
}
