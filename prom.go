package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)
##
type PrometheusService struct {
	testCounter    prometheus.Counter
	testCounterVec *prometheus.CounterVec
	testGauge      prometheus.Gauge
	testHistogram  prometheus.Histogram
	testSummary    prometheus.Summary
}

func NewPrometheusService(reg prometheus.Registerer) *PrometheusService {
	// 使用promauto包中的函数来实现自动注册
	// create a Factory once to be used multiple times
	factory := promauto.With(reg)

	// Counter
	// 不加label
	opsProcessed := factory.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
	// 加label, 扩展性好，比如对于http请求中多个不同路径的数据统计
	opsProcessedLabel := factory.NewCounterVec(prometheus.CounterOpts{
		Name: "myapp_processed_ops_label_total",
		Help: "The total number of processed events",
	}, []string{"myapp", "path"})

	// Gauge
	cpuTemp := factory.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	})

	// Histogram
	histogram := factory.NewHistogram(prometheus.HistogramOpts{
		Name:    "integer_histogram_numbers",
		Help:    "A histogram of integer numbers.",
		Buckets: prometheus.LinearBuckets(1, 1, 10),
	})

	// Summary
	integerSummary := factory.NewSummary(
		prometheus.SummaryOpts{
			Name:       "integer_summary_numbers",
			Help:       "Summary of integer numbers",
			Objectives: map[float64]float64{0.2: 0.001, 0.5: 0.001, 0.8: 0.001, 0.99: 0.001}, // 容忍的误差根据需求自定义
		},
	)

	return &PrometheusService{
		testCounter:    opsProcessed,
		testCounterVec: opsProcessedLabel,
		testGauge:      cpuTemp,
		testHistogram:  histogram,
		testSummary:    integerSummary,
	}
}
