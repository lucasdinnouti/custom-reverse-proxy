package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"proxy/selectors"
	"proxy/targets"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	routeSelector Selector
	targetProxy   map[string]*httputil.ReverseProxy

	requestDurations = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "proxy_request_latency",
		Help:    "Latency of requests received from runner",
		Buckets: []float64{0.001, 0.005, 0.010, 0.100, 0.500, 1, 5, 10}})
)

type Selector interface {
	Select(*http.Request) (string, error)
	Destroy()
}

func route(w http.ResponseWriter, r *http.Request) {
	target, err := routeSelector.Select(r)

	if err != nil {
		log.Fatal(err)
		return
	}

	if fn, ok := targetProxy[target]; ok {
		log.Println("target: ", target)

		before := time.Now()
		timer := prometheus.NewTimer(requestDurations)
		fn.ServeHTTP(w, r)
		timer.ObserveDuration()
		log.Println("Time elapsed", time.Since(before))
		log.Println("")

		return
	}

	w.Write([]byte("403: Host forbidden " + target))
}

func main() {
	prometheus.MustRegister(requestDurations)
	http.Handle("/metrics", promhttp.Handler())

	hosts := []string{"a", "b", "c"}
	weights := []int{2, 1, 1}
	types := map[string][]string{"image": []string{"c"}}
	nodeTypes := map[string]string{"a": "large-cpu", "b": "medium-cpu", "c": "medium-gpu"}

	switch algorithm := os.Getenv("ALGORITHM"); algorithm {
	case "round_robin":
		routeSelector = selectors.NewRoundRobin(hosts)
	case "weighted_round_robin":
		routeSelector = selectors.NewWeightedRoundRobin(hosts, weights)
	case "metadata":
		routeSelector = selectors.NewMetadata(hosts, types)
	case "machine_learning":
		routeSelector = selectors.NewMachineLearning(hosts, nodeTypes)
	default:
		routeSelector = selectors.NewRoundRobin(hosts)
	}

	targetProxy = targets.Build(hosts)
	http.HandleFunc("/message", route)

	http.ListenAndServe(":8082", nil)
	routeSelector.Destroy()
}
