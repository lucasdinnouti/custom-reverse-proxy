package main

import (
    "fmt"
    "github.com/owulveryck/onnx-go"
    "github.com/owulveryck/onnx-go/backend/x/gorgonnx"
    "gorgonia.org/tensor"
    "os"
    "log"
)

func main() {

    // Create a backend receiver
	backend := gorgonnx.NewGraph()
	// Create a model and set the execution backend
	model := onnx.NewModel(backend)

	// read the onnx model
	b, _ := os.ReadFile("test_2.onnx")
	// Decode it into the model
	err := model.UnmarshalBinary(b)
	if err != nil {
		log.Fatal(err)
	}

    // Create a tensor with the input
    inputData := []float32{0, 1, 0.319191, 0.669730, 0.635886, 3.475952, 5.505371, 7.519531}
    inputTensor := tensor.New(tensor.WithShape(1, len(inputData)), tensor.Of(tensor.Float32), tensor.WithBacking(inputData))

	// Set the first input, the number depends of the model
	model.SetInput(0, inputTensor)

	err = backend.Run()
	if err != nil {
		log.Fatal(err)
	}
    
	output, _ := model.GetOutputTensors()
	fmt.Println(output[0])

}