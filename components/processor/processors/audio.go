package processors

type Audio struct {
}

func NewAudio() *Audio {
	return &Audio{}
}

func (a *Audio) Process(content string) string {
	
	WasteTime(200, 220)

	return content
}
