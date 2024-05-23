package targets

import (
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
)

func Build(hosts []string) map[string]*httputil.ReverseProxy {
	targetProxy := map[string]*httputil.ReverseProxy{}

	for _, h := range hosts {

		dns := fmt.Sprintf("http://processor-%s.default.svc.cluster.local:8083", h)
		remoteUrl, err := url.Parse(dns)
		if err != nil {
			log.Println("target parse fail:", err)
			return targetProxy
		}
		targetProxy[h] = httputil.NewSingleHostReverseProxy(remoteUrl)
	}

	return targetProxy
}
