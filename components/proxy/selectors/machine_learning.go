package selectors

import "net/http"

type MachineLearning struct {
}

func NewMachineLearning() *MachineLearning {
	return &MachineLearning{}
}

func (r *MachineLearning) Select(request *http.Request) (string, error) {
	return "a", nil
}
