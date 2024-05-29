package processors

type Default struct {
}

func NewDefault() *Default {
	return &Default{}
}

func (d *Default) Process(content string) string {

	WasteTime(10, 100)

	return content
}
