package processors

import (
	"os"
	"strings"
)

type Image struct {
}

func NewImage() *Image {
	return &Image{}
}

func (i *Image) Process(content string) string {
	instance := os.Getenv("INSTANCE_TYPE")

	if strings.Contains(instance, "gpu") {
		WasteTime(10, 100)
	} else {
		WasteTime(100, 1000)
	}

	return content
}
