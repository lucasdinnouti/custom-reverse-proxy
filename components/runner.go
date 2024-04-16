package main

import (
    "fmt"
    "log"
    "net/http"
    "strconv"
    "sync"
)

var hits map[string]int
var mutex = &sync.Mutex{}

func route(w http.ResponseWriter, r *http.Request) {
    mutex.Lock()
    hits["instance"]++
	
	// TODO: request to proxy, and increment called processor instance 

    fmt.Fprintf(w, strconv.Itoa(hits["instance"]))
    mutex.Unlock()
}

func main() {
    http.HandleFunc("/route", route)

    http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hi")
    })

    log.Fatal(http.ListenAndServe(":8081", nil))
}