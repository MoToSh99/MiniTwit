package main

import (
        "net/http"
        "time"
			
        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
        go func() {
                for {
						CPU_GAUGE.Inc()
						REPONSE_COUNTER.Inc()
						REQ_DURATION_SUMMARY.Observe()
                        time.Sleep(2 * time.Second)
                }
        }()
}

var (
		CPU_GAUGE  = promauto.NewGauge(prometheus.GaugeOpts{
                Name: "minitwit_cpu_load_percent",
                Help: "Current load of the CPU in percent.",
		})
		
		REPONSE_COUNTER = promauto.NewCounter(prometheus.CounterOpts{
			Name: "myapp_processed_ops_total",
			Help: "The count of HTTP responses sent.",
		})

		REQ_DURATION_SUMMARY = promauto.NewHistogram(prometheus.HistogramOpts{
			Name: "minitwit_request_duration_milliseconds",
			Help: "Request duration distribution.",
		})
)

func main() {
        recordMetrics()

        http.Handle("/metrics", promhttp.Handler())
        http.ListenAndServe(":2112", nil)
}