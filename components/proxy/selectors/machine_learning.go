package selectors

import (
	"log"
	"net/http"

	ort "github.com/yalue/onnxruntime_go"
)

var MODEL_FILE = "model.onnx"

type MachineLearning struct {
	Session    *ort.DynamicAdvancedSession
	Hosts      []string
	HostsCount int
	Types      map[string]string
	Translator map[string]float32
}

func NewMachineLearning(hosts []string, types map[string]string) *MachineLearning {
	ort.SetSharedLibraryPath("./libonnxruntime.so")

	err := ort.InitializeEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	session, err := ort.NewDynamicAdvancedSession(
		MODEL_FILE,
		[]string{"float_input"},
		[]string{"elapsed"},
		nil)

	if err != nil {
		log.Fatal(err)
	}

	translator := map[string]float32{
		"a":          0.0,
		"b":          1.0,
		"c":          2.0,
		"text":       0.0,
		"image":      1.0,
		"audio":      2.0,
		"large-cpu":  0.0,
		"medium-cpu": 1.0,
		"medium-gpu": 2.0,
	}

	return &MachineLearning{
		Session:    session,
		Hosts:      hosts,
		HostsCount: len(hosts),
		Types:      types,
		Translator: translator,
	}
}

func (r *MachineLearning) Destroy() {
	ort.DestroyEnvironment()
	r.Session.Destroy()
}

func (r *MachineLearning) buildInputTensor(request *http.Request) (*ort.Tensor[float32], error) {
	inputData := []float32{}

	messageType := request.Header.Get("X-Message-Type")
	processorAcpu := float32(4.346822)
	processorBcpu := float32(6.074483)
	processorCcpu := float32(14.455157)
	processorAmem := float32(10.81543)
	processorBmem := float32(10.229492)
	processorCmem := float32(18.200684)

	for i := 0; i < r.HostsCount; i++ {
		instanceId := r.Hosts[i]
		instanceType := r.Types[instanceId]

		inputData = append(inputData, r.Translator[messageType])
		inputData = append(inputData, r.Translator[instanceId])
		inputData = append(inputData, r.Translator[instanceType])
		inputData = append(inputData, processorAcpu)
		inputData = append(inputData, processorBcpu)
		inputData = append(inputData, processorCcpu)
		inputData = append(inputData, processorAmem)
		inputData = append(inputData, processorBmem)
		inputData = append(inputData, processorCmem)
	}

	inputShape := ort.NewShape(int64(r.HostsCount), 9)
	inputTensor, err := ort.NewTensor(inputShape, inputData)
	if err != nil {
		return nil, err
	}

	return inputTensor, nil
}

func (r *MachineLearning) buildOutputTensor(request *http.Request) (*ort.Tensor[float32], error) {
	outputShape := ort.NewShape(int64(r.HostsCount), 1)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		return nil, err
	}

	return outputTensor, nil
}

func (r *MachineLearning) getBetterResult(results []float32) string {
	min := results[0]
	minIndex := 0
	for i := 1; i < len(results); i++ {
		if results[i] < min {
			min = results[i]
			minIndex = i
		}
	}

	log.Println(min)
	return r.Hosts[minIndex]
}

func (r *MachineLearning) Select(request *http.Request) (string, error) {
	inputTensor, err := r.buildInputTensor(request)
	defer inputTensor.Destroy()
	if err != nil {
		return "", err
	}

	outputTensor, err := r.buildOutputTensor(request)
	defer outputTensor.Destroy()
	if err != nil {
		return "", err
	}

	err = r.Session.Run(
		[]ort.ArbitraryTensor{inputTensor},
		[]ort.ArbitraryTensor{outputTensor})

	host := r.getBetterResult(outputTensor.GetData())

	return host, nil
}
