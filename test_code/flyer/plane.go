package flyer

type Wing interface {
	Material() string
}

type AluminumWing struct {
}

func (*AluminumWing) Material() string {
	return "aluminum"
}

type Plane struct {
	Wing Wing `inject:""`
}
