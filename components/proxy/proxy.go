package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"proxy/selectors"
)

var (
	routeSelector Selector
	targetProxy   map[string]*httputil.ReverseProxy = map[string]*httputil.ReverseProxy{}
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

	remoteUrl, err := url.Parse("http://processor-a.default.svc.cluster.local:8083")
	if err != nil {
		log.Println("target parse fail:", err)
		return
	}
	targetProxy["a"] = httputil.NewSingleHostReverseProxy(remoteUrl)

	remoteUrl, err = url.Parse("http://processor-b.default.svc.cluster.local:8083")
	if err != nil {
		log.Println("target parse fail:", err)
		return
	}
	targetProxy["b"] = httputil.NewSingleHostReverseProxy(remoteUrl)

	http.HandleFunc("/message", route)

	http.ListenAndServe(":8082", nil)
}
