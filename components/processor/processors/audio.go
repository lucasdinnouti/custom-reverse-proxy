package processors

import (
	"time"
)

type Audio struct {
}

func NewAudio() *Audio {
	return &Audio{}
}

func (a *Audio) Process(content string) string {

	time.Sleep(500 * time.Millisecond)
	return content
}
