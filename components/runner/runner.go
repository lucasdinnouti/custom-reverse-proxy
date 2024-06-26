package main

import (
	"log"
	"net/http"

	"runner/loadtest"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestDurations = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "full_request_latency",
		Help:    "Latency of requests to processor",
		Buckets: []float64{0.001, 0.002, 0.003, 0.004, 0.005, 0.006, 0.007, 0.008, 0.009, 0.010, 0.100, 0.500, 1, 2, 3, 4, 5, 6, 8, 9, 10, 15, 20, 25, 30}},
		[]string{"routed_to"},
	)

	requestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_status",
		Help: "Number of requests by status"},
		[]string{"status"},
	)
)

func main() {

	prometheus.MustRegister(requestDurations)
	prometheus.MustRegister(requestCounter)

	http.Handle("/metrics", promhttp.Handler())

	go http.ListenAndServe(":8081", nil)

	time.Sleep(30 * time.Second)

	log.Println("Starting Runner!")

	loadtest.LoadTestCase("testcase.txt")
	loadtest.RunTestCase(requestDurations, requestCounter)

	time.Sleep(12 * time.Hour)
}
