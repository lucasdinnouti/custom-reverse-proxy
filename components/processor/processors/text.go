package processors

import (
	"time"
)

type Text struct {
}

func NewText() *Text {
	return &Text{}
}

func (t *Text) Process(content string) string {
	time.Sleep(10 * time.Microsecond)

	return content
}
