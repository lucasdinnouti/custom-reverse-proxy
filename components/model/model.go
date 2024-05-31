package main

import (
	"fmt"
	"log"

	ort "github.com/yalue/onnxruntime_go"
)

//func main() {
//
//    // Create a backend receiver
//	backend := gorgonnx.NewGraph()
//	// Create a model and set the execution backend
//	model := onnx.NewModel(backend)
//
//	// read the onnx model
//	b, _ := os.ReadFile("test_2.onnx")
//	// Decode it into the model
//	err := model.UnmarshalBinary(b)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//    // Create a tensor with the input
//    inputData := []float32{0, 1, 0.319191, 0.669730, 0.635886, 3.475952, 5.505371, 7.519531}
//    inputTensor := tensor.New(tensor.WithShape(1, len(inputData)), tensor.Of(tensor.Float32), tensor.WithBacking(inputData))
//
//	// Set the first input, the number depends of the model
//	model.SetInput(0, inputTensor)
//
//	err = backend.Run()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	output, _ := model.GetOutputTensors()
//	fmt.Println(output[0])
//
//}

func main() {
	ort.SetSharedLibraryPath("./libonnxruntime.so")

	err := ort.InitializeEnvironment()
	defer ort.DestroyEnvironment()

	if err != nil {
		log.Fatal(err)
	}

	inputData := []float32{0.0, 1.0, 1.0, 0.319191, 0.669730, 0.635886, 3.475952, 5.505371, 7.519531}
	inputShape := ort.NewShape(1, 9)
	inputTensor, err := ort.NewTensor(inputShape, inputData)
	defer inputTensor.Destroy()

	if err != nil {
		log.Fatal(err)
	}

	outputShape := ort.NewShape(1, 1)
	outputTensor, _ := ort.NewEmptyTensor[float32](outputShape)
	defer outputTensor.Destroy()

	session, err := ort.NewAdvancedSession("./test_3.onnx",
		[]string{"float_input"}, []string{"elapsed"},
		[]ort.ArbitraryTensor{inputTensor}, []ort.ArbitraryTensor{outputTensor}, nil)
	defer session.Destroy()

	if err != nil {
		log.Fatal(err)
	}
	// Calling Run() will run the network, reading the current contents of the
	// input tensors and modifying the contents of the output tensors.
	err = session.Run()

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
