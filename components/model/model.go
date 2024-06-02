package main

import (
	"fmt"
	"log"

	ort "github.com/yalue/onnxruntime_go"
)

func main() {
	ort.SetSharedLibraryPath("./libonnxruntime.so")

	err := ort.InitializeEnvironment()
	defer ort.DestroyEnvironment()

	if err != nil {
		log.Fatal(err)
	}

	inputData := []float32{0.0, 1.0, 1.0, 0.319191, 0.669730, 0.635886, 3.475952, 5.505371, 7.519531, 0.0, 2.0, 2.0, 0.319191, 0.669730, 0.635886, 3.475952, 5.505371, 7.519531}
	inputShape := ort.NewShape(2, 9)
	inputTensor, err := ort.NewTensor(inputShape, inputData)
	defer inputTensor.Destroy()

	if err != nil {
		log.Fatal(err)
	}

	outputShape := ort.NewShape(2, 1)
	outputTensor, _ := ort.NewEmptyTensor[float32](outputShape)
	defer outputTensor.Destroy()

	session, err := ort.NewDynamicAdvancedSession("./test_3.onnx",
		[]string{"float_input"}, []string{"elapsed"}, nil)

	defer session.Destroy()

	if err != nil {
		log.Fatal(err)
	}

	err = session.Run([]ort.ArbitraryTensor{inputTensor}, []ort.ArbitraryTensor{outputTensor})

	if err != nil {
		log.Fatal(err)
	}

	// Get a slice view of the output tensor's data.
	outputData := outputTensor.GetData()

	fmt.Println(outputData)
	// If you want to run the network on a different input, all you need to do
	// is modify the input tensor data (available via inputTensor.GetData())
	// and call Run() again.

	// ...
}
