package cardstack

type Card struct {
	Name     string `json:"name"`
	ImageSrc string `json:"image_src"`
}

type Stack struct {
	Cards map[int]Card `json:"stack"`
}

type CardDto struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
