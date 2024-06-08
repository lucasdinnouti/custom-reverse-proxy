package processors

import (
	"os"
	"strings"
	"time"
)

type Image struct {
}

func NewImage() *Image {
	return &Image{}
}

func (i *Image) Process(content string) string {
	instance := os.Getenv("INSTANCE_TYPE")

	if strings.Contains(instance, "gpu") {
		time.Sleep(100 * time.Microsecond)
	} else {
		time.Sleep(1000 * time.Microsecond)
	}

	return content
}
