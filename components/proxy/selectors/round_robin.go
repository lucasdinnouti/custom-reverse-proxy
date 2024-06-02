package selectors

import "net/http"

type RoundRobin struct {
	Counter    int
	Hosts      []string
	HostsCount int
}

func NewRoundRobin(hosts []string) *RoundRobin {
	return &RoundRobin{
		Counter:    0,
		Hosts:      hosts,
		HostsCount: len(hosts),
	}
}

func (r *RoundRobin) Destroy() {}

func (r *RoundRobin) Select(request *http.Request) (string, error) {
	r.Counter++

	return r.Hosts[r.Counter%r.HostsCount], nil
}
