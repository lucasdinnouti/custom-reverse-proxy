package selectors

type RoundRobin struct {
	Counter int
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{
		Counter: 0,
	}
}

func (r *RoundRobin) Select() string {
	r.Counter++

	if r.Counter%2 == 0 {
		return "a"
	}

	return "b"
}
