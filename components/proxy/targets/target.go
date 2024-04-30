package targets

import (
	"log"
	"net/http/httputil"
	"net/url"
)

func Get() map[string]*httputil.ReverseProxy {
	targetProxy := map[string]*httputil.ReverseProxy{}

	remoteUrl, err := url.Parse("http://processor-a.default.svc.cluster.local:8083")
	if err != nil {
		log.Println("target parse fail:", err)
		return targetProxy
	}
	targetProxy["a"] = httputil.NewSingleHostReverseProxy(remoteUrl)

	remoteUrl, err = url.Parse("http://processor-b.default.svc.cluster.local:8083")
	if err != nil {
		log.Println("target parse fail:", err)
		return targetProxy
	}
	targetProxy["b"] = httputil.NewSingleHostReverseProxy(remoteUrl)

	return targetProxy
}
