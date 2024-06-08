package processors

import (
	"time"
)

type Default struct {
}

func NewDefault() *Default {
	return &Default{}
}

func (d *Default) Process(content string) string {

	time.Sleep(10 * time.Microsecond)

	return content
}
