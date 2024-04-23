package main

import (
	"fmt"
	"log"
	"net/http"
)

func echoString(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "a")
}

func main() {

	http.HandleFunc("/", echoString)

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Test request")
		fmt.Fprintf(w, "Test")
	})

	log.Fatal(http.ListenAndServe(":8083", nil))
}
