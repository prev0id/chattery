package icon

type iconType int

const (
	defaultSize        = 24
	defaultFill        = "none"
	defaultStroke      = "currentColor"
	defaultStrokeWidth = "2"
)

type Props struct {
	Size        int
	Color       string
	Fill        string
	Stroke      string
	StrokeWidth string
	Class       string
}

func (i *Props) GetSize() int {
	if i.Size <= 0 {
		return defaultSize
	}
	return i.Size
}

func (i *Props) GetFill() string {
	if i.Fill == "" {
		return defaultFill
	}
	return i.Fill
}

func (i *Props) GetStroke() string {
	stroke := i.Stroke
	if stroke == "" {
		stroke = i.Color
	}
	if stroke == "" {
		stroke = defaultStroke
	}
	return stroke
}

func (i *Props) GetStrokeWidth() string {
	if i.StrokeWidth == "" {
		return defaultStrokeWidth
	}
	return i.StrokeWidth
}
