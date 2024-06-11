package processors

type Text struct {
}

func NewText() *Text {
	return &Text{}
}

func (t *Text) Process(content string) string {
	WasteTime(10, 100)

	return content
}
