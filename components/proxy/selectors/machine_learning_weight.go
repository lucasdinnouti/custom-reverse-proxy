package selectors

import (
	"errors"
	"net/http"

	"proxy/jobs"
)

type MachineLearningWeight struct {
	Counter     int
	Hosts       []string
	HostsCount  int
	WeightCache *jobs.WeightCache
}

func NewMachineLearningWeight(hosts []string, types map[string]string, initialWeights []int) *MachineLearningWeight {
	weightCache := jobs.NewWeightCache(hosts, types, initialWeights)
	go weightCache.Run()

	return &MachineLearningWeight{
		Counter:     0,
		Hosts:       hosts,
		HostsCount:  len(hosts),
		WeightCache: weightCache,
	}
}

func (r *MachineLearningWeight) Destroy() {
	r.WeightCache.Destroy()
}

func (r *MachineLearningWeight) Select(request *http.Request) (string, error) {
	r.Counter++

	messageType := request.Header.Get("X-Message-Type")
	r.WeightCache.IncType(messageType)

	weightsSum := 0
	for _, v := range r.WeightCache.Weights {
		weightsSum += v
	}

	acc := 0
	for i := 0; i < r.HostsCount; i++ {
		if r.Counter%weightsSum < r.WeightCache.Weights[i]+acc {
			return r.Hosts[i], nil
		}
		acc += r.WeightCache.Weights[i]
	}

	return "", errors.New("failed to select host")
}
