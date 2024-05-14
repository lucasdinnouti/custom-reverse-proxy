package selectors

import "errors"

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

func (r *RoundRobin) Select() (string, error) {
	r.Counter++

	for i := 0; i < r.HostsCount; i++ {
		if r.Counter%r.HostsCount == i {
			return r.Hosts[i], nil
		}
	}

	return "", errors.New("faled to select host")
}
