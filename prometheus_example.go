package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Registry = prometheus.DefaultRegisterer //可以全局调用注册

	// 非自动注册，见init函数
	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors.",
		},
		[]string{"device"},
	)
)

func init() {
	// Metrics have to be registered to be exposed:
	Registry.MustRegister(hdFailures)
}

func main() {
	p := NewPrometheusService(Registry)

	// mock 程序中调用Metrics
	recordMetrics(p)
	gaugeMetrics(p)
	histogramSummary(p)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2112", nil))
}

func recordMetrics(p *PrometheusService) {
	go func() {
		for {
			p.testCounter.Inc()
			p.testCounterVec.WithLabelValues("app1", "/bar").Inc()
			p.testCounterVec.WithLabelValues("app1", "/foo").Inc()

			hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

func gaugeMetrics(p *PrometheusService) {
	go func() {
		p.testGauge.Set(65.3)
		rand.Seed(time.Now().Unix())
		for {
			p.testGauge.Add(float64(rand.Intn(11) - 5)) // 可增可减
			time.Sleep(2 * time.Second)

		}
	}()
}

func histogramSummary(p *PrometheusService) {
	go func() {
		for {
			for i := 1.0; i <= 10; i++ {
				p.testHistogram.Observe(i)
				p.testSummary.Observe(i)
				time.Sleep(2 * time.Second)
			}
		}
	}()
}
