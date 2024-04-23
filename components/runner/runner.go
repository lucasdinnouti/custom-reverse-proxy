package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

var hits map[string]int
var mutex = &sync.Mutex{}

func request() {
	mutex.Lock()
	hits["instance"]++

	log.Println("Requesting processor...")
	_, err := http.Get("http://proxy.default.svc.cluster.local:8082/test")
	if err != nil {
		log.Fatalln(err)
	}

	mutex.Unlock()
}

func main() {
	hits = map[string]int{
		"instance": 0,
	}

	for {
		request()
		time.Sleep(5 * time.Second)
	}
}
