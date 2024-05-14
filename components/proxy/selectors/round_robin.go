package selectors

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

	return r.Hosts[r.Counter%r.HostsCount], nil
}
