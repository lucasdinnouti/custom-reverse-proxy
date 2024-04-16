package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
    counter int
	targetProxy map[string]*httputil.ReverseProxy = map[string]*httputil.ReverseProxy{}
)

func select_host() string {
    counter++

    if counter%2 == 0 {
        return "a" 
    } else {
        return "b" 
    }
}

func route(w http.ResponseWriter, r *http.Request) {
	target := select_host()

	if fn, ok := targetProxy[target]; ok {
        log.Println("target: ", target)
		
        fn.ServeHTTP(w, r)
        
		return
	}
    
	w.Write([]byte("403: Host forbidden " + target))
}

func main() {

    remoteUrl, err := url.Parse("http://localhost:8083")
    if err != nil {
        log.Println("target parse fail:", err)
        return
    }
    
    targetProxy["a"] = httputil.NewSingleHostReverseProxy(remoteUrl)

    remoteUrl, err = url.Parse("http://localhost:8084")
    if err != nil {
        log.Println("target parse fail:", err)
        return
    }
    
    targetProxy["b"] = httputil.NewSingleHostReverseProxy(remoteUrl)

	http.HandleFunc("/test", route)

    http.ListenAndServe(":8082", nil)
}