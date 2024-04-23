package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Message struct {
	datetime string `json:"datetime"`
	content string `json:"content"`
}

func echoString(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "a")
}

func parseMessage(w http.ResponseWriter, r *http.Request) {
	m := &Message{}

	err := json.NewDecoder(r.Body).Decode(m)

	if (err != nil) {
		log.Fatalln("Error Parsing Request Body", r.Body)
	}

	log.Println(r.Body)
	log.Println(m)

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
