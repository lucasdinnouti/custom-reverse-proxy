package selectors

type Metadata struct {
}

func NewMetadata() Metadata {
	return Metadata{}
}

func (r Metadata) Select() string {
	return "a"
}
