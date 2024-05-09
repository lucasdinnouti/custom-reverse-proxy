package selectors

import (
	"log"
)

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

	log.Println(r.Counter)

	if r.Counter%2 == 0 {
		return "a"
	}

	return "b"
}
