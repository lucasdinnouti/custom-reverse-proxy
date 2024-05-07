package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
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

var rateLimiter = time.Tick(5 * time.Millisecond)

func echoString(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "a")
}

func parseMessage(w http.ResponseWriter, r *http.Request) {
	m := &Message{}

	err := json.NewDecoder(r.Body).Decode(m)

	if err != nil {
		log.Fatalln("Error Parsing Request Body", r.Body)
	}

	log.Println(r.Body)
	log.Println(m)

	<-rateLimiter

	fmt.Fprintf(w, "OK")
}

func main() {

	http.HandleFunc("/", echoString)

	http.HandleFunc("/message", parseMessage)

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Test request")
		fmt.Fprintf(w, "Test")
	})

	log.Fatal(http.ListenAndServe(":8083", nil))
}
