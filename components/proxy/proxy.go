package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"proxy/selectors"
	"proxy/targets"
)

var (
	routeSelector Selector
	targetProxy   map[string]*httputil.ReverseProxy
)

type Selector interface {
	Select() string
}

func route(w http.ResponseWriter, r *http.Request) {
	target := routeSelector.Select()

	if fn, ok := targetProxy[target]; ok {
		log.Println("target: ", target)

		fn.ServeHTTP(w, r)

		return
	}

	w.Write([]byte("403: Host forbidden " + target))
}

func main() {

	switch algorithm := os.Getenv("ALGORITHM"); algorithm {
	case "round_robin":
		routeSelector = selectors.NewRoundRobin()
	case "weighted_round_robin":
		routeSelector = selectors.NewWeightedRoundRobin()
	case "metadata":
		routeSelector = selectors.NewMetadata()
	case "machine_learning":
		routeSelector = selectors.NewMachineLearning()
	default:
		routeSelector = selectors.NewRoundRobin()
	}

	targetProxy = targets.Get()
	http.HandleFunc("/message", route)

	http.ListenAndServe(":8082", nil)
}
