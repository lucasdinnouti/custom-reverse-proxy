package selectors

type WeightedRoundRobin struct {
}

func NewWeightedRoundRobin() *WeightedRoundRobin {
	return &WeightedRoundRobin{}
}

func (r *WeightedRoundRobin) Select() string {
	return "a"
}
