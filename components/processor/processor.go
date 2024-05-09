package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ContentType uint8
const (
  Text ContentType = iota
  Image
  Audio
  Unknown
)

type Message struct {
	Datetime string 	 `json:"datetime"`
	Content  string 	 `json:"content"`
	Type 	 ContentType `json:"type"`
}

var counter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "request_count",
	Help: "Number of requests received from runner"})

func echoString(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "a")
}

func parseMessage(w http.ResponseWriter, r *http.Request) {
	counter.Inc()

	m := &Message{}

	err := json.NewDecoder(r.Body).Decode(m)

	if err != nil {
		log.Fatalln("Error Parsing Request Body", r.Body)
	}

	log.Println(r.Body)
	log.Println(m)

	time.Sleep(1 * time.Second)
	fmt.Fprintf(w, "OK")
}

func main() {

	prometheus.MustRegister(counter)

	http.HandleFunc("/", echoString)

	http.HandleFunc("/message", parseMessage)

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Test request")
		fmt.Fprintf(w, "Test")
	})

	log.Fatal(http.ListenAndServe(":8083", nil))
}
