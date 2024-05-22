package main

import (
	"processor/processors"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	
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

type Processor interface {
	Process(content string) string
}

var (
	textProcessor = processors.NewText()
	imageProcessor = processors.NewImage()
	audioProcessor = processors.NewAudio()
	defaultProcessor = processors.NewDefault()

	counter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "request_count",
		Help: "Number of requests received from runner"})
)

func processMessage(message *Message) {
	switch message.Type {
	case Text:
		textProcessor.Process(message.Content)
	case Image:
		imageProcessor.Process(message.Content)
	case Audio:
		audioProcessor.Process(message.Content)
	default:
		defaultProcessor.Process(message.Content)
	}
}

func handleMessage(w http.ResponseWriter, r *http.Request) {
	counter.Inc()

	message := &Message{}

	err := json.NewDecoder(r.Body).Decode(message)

	if err != nil {
		log.Fatalln("Error Parsing Request Body", r.Body)
	}

	log.Println(message)
	processMessage(message)

	fmt.Fprintf(w, os.Getenv("INSTANCE_TYPE"))
}

func main() {

	prometheus.MustRegister(counter)

	http.HandleFunc("/message", handleMessage)

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Test request")
		fmt.Fprintf(w, "Test")
	})

	log.Fatal(http.ListenAndServe(":8083", nil))
}
