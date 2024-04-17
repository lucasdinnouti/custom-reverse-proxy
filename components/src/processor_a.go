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

    http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hi")
    })

    log.Fatal(http.ListenAndServe(":8083", nil))
}