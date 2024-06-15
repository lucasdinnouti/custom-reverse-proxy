package jobs

import (
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	ort "github.com/yalue/onnxruntime_go"
)

var MODEL_FILE = "model-weight.onnx"

type WeightCache struct {
	Weights         []int
	IntervalSeconds int
	TypeCounter     map[string]int
	Translator      map[string]float32
	Session         *ort.DynamicAdvancedSession
	Hosts           []string
	HostsCount      int
	Types           map[string]string
	PromCache       *PromCache
}

func NewWeightCache(hosts []string, types map[string]string, initialWeights []int) *WeightCache {
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

	promCache := NewPromCache(hosts)
	go promCache.Run()

	intervalSeconds, err := strconv.Atoi(os.Getenv("WEIGHT_INTERVAL"))
	if err != nil {
		log.Fatal(err)
	}

	return &WeightCache{
		Weights:         initialWeights,
		IntervalSeconds: intervalSeconds,
		TypeCounter: map[string]int{
			"image": 0,
			"text":  0,
			"audio": 0},
		PromCache:  promCache,
		Translator: translator,
		Session:    session,
		Hosts:      hosts,
		HostsCount: len(hosts),
		Types:      types,
	}
}

func (r *WeightCache) Destroy() {
	ort.DestroyEnvironment()
	r.Session.Destroy()
}

func (r *WeightCache) IncType(messageType string) {
	r.TypeCounter[messageType]++
}

func (r *WeightCache) getBiggestType() string {
	max := -1
	maxType := ""

	for k, v := range r.TypeCounter {
		if v > max || max == -1 {
			max = v
			maxType = k
		}
	}

	return maxType
}

func (r *WeightCache) buildInputTensor() (*ort.Tensor[float32], error) {
	inputData := []float32{}

	messageType := r.getBiggestType()
	for _, instanceId := range r.Hosts {
		instanceType := r.Types[instanceId]

		inputData = append(inputData, r.Translator[messageType])
		inputData = append(inputData, r.Translator[instanceId])
		inputData = append(inputData, r.Translator[instanceType])

		for _, instanceId2 := range r.Hosts {
			inputData = append(inputData, r.PromCache.CpuUsage[instanceId2])
		}

		for _, instanceId2 := range r.Hosts {
			inputData = append(inputData, r.PromCache.MemoryUsage[instanceId2])
		}
	}

	inputShape := ort.NewShape(int64(r.HostsCount), 9)
	inputTensor, err := ort.NewTensor(inputShape, inputData)
	if err != nil {
		return nil, err
	}

	return inputTensor, nil
}

func (r *WeightCache) buildOutputTensor() (*ort.Tensor[float32], error) {
	outputShape := ort.NewShape(int64(r.HostsCount), 1)
	outputTensor, err := ort.NewEmptyTensor[float32](outputShape)
	if err != nil {
		return nil, err
	}

	return outputTensor, nil
}

func (r *WeightCache) setWeights(results []float32) {
	latencyToHost := map[int]int{}
	latencyList := []int{}
	for k, v := range results {
		latencyToHost[int(v)] = k
		latencyList = append(latencyList, int(v))
	}

	sort.Ints(latencyList)

	weights := []int{3, 2, 1}
	for k, v := range latencyList {
		i := latencyToHost[v]
		r.Weights[i] = weights[k]
	}

	log.Println(results)
	log.Println(r.Weights)
}

func (r *WeightCache) Run() {
	for {
		log.Println("Starting loading weights...")

		inputTensor, err := r.buildInputTensor()
		defer inputTensor.Destroy()
		if err != nil {
			log.Fatal(err)
		}

		outputTensor, err := r.buildOutputTensor()
		defer outputTensor.Destroy()
		if err != nil {
			log.Fatal(err)
		}

		err = r.Session.Run(
			[]ort.ArbitraryTensor{inputTensor},
			[]ort.ArbitraryTensor{outputTensor})

		time.Sleep(time.Duration(r.IntervalSeconds) * time.Second)

		r.setWeights(outputTensor.GetData())
		log.Println("WeightCache sleep...")
	}
}
