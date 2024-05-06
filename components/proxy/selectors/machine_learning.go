package selectors

type MachineLearning struct {
}

func NewMachineLearning() *MachineLearning {
	return &MachineLearning{}
}

func (r *MachineLearning) Select() string {
	return "a"
}
