package views

type Renderer struct {
	asset func(string) string
}

func NewRenderer(asset func(string) string) Renderer {
	return Renderer{asset: asset}
}