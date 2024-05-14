package processors

import (
	"time"
)

type Image struct {
}

func NewImage() *Image {
	return &Image{}
}

func (i *Image) Process(content string) string {

	time.Sleep(1000 * time.Millisecond)
	return content
}
