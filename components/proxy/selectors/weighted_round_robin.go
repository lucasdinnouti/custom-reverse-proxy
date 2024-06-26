package selectors

import (
	"errors"
	"net/http"
)

type WeightedRoundRobin struct {
	Counter    int
	Hosts      []string
	HostsCount int
	Weights    []int
	WeightsSum int
}

func NewWeightedRoundRobin(hosts []string, weights []int) *WeightedRoundRobin {
	weightsSum := 0
	for _, v := range weights {
		weightsSum += v
	}

	return &WeightedRoundRobin{
		Counter:    0,
		Hosts:      hosts,
		HostsCount: len(hosts),
		Weights:    weights,
		WeightsSum: weightsSum,
	}
}

func (r *WeightedRoundRobin) Destroy() {}

func (r *WeightedRoundRobin) Select(request *http.Request) (string, error) {
	r.Counter++

	acc := 0
	for i := 0; i < r.HostsCount; i++ {
		if r.Counter%r.WeightsSum < r.Weights[i]+acc {
			return r.Hosts[i], nil
		}
		acc += r.Weights[i]
	}

	return "", errors.New("failed to select host")
}
