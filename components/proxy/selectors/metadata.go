package selectors

import (
	"log"
	"math/rand"
	"net/http"
)

type Metadata struct {
	Hosts []string
	Types map[string][]string
}

func NewMetadata(hosts []string, types map[string][]string) *Metadata {
	return &Metadata{
		Hosts: hosts,
		Types: types,
	}
}

func randRange(min, max int) int {
	return rand.Intn(max-min) + min
}

func (r *Metadata) Destroy() {}

func (r *Metadata) Select(request *http.Request) (string, error) {
	messageType := request.Header.Get("X-Message-Type")
	specializedHosts := r.Types[messageType]

	log.Println(messageType)
	if specializedHosts != nil {
		specializedHostsCount := len(specializedHosts)
		if specializedHostsCount > 1 {
			randIndex := randRange(0, specializedHostsCount)
			return specializedHosts[randIndex], nil
		}

		return specializedHosts[0], nil
	}

	randIndex := randRange(0, len(r.Hosts))
	return r.Hosts[randIndex], nil
}
