package cardstack

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
)

type Card struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	ImageSrc string `json:"image_src"`
}

/* ToDo: Use Map to avoid getImageURL Function */
type Stack struct {
	Cards []Card `json:"stack"`
}

type CardDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (stack *Stack) GetImageURL(id string) string {
	for _, card := range stack.Cards {
		if card.Id == id {
			return card.ImageSrc
		}
	}
	return ""
}

func (stack *Stack) Draw3() []CardDto {
	stack.shuffle()
	var dto []CardDto
	for i := 0; i < 3; i++ {
		dto = append(dto, CardDto{stack.Cards[i].Id, stack.Cards[i].Name})
	}
	return dto
}

func (stack *Stack) shuffle() {
	// very important information: fisher-yates algo
	for i := len(stack.Cards) - 1; i > 0; i-- {
		randPos := rand.Intn(i + 1)
		stack.Cards[i], stack.Cards[randPos] = stack.Cards[randPos], stack.Cards[i]
	}
}

func NewStack(path string) *Stack {
	stack := &Stack{}
	stackFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer stackFile.Close()

	err = json.NewDecoder(stackFile).Decode(stack)
	if err != nil {
		fmt.Println(err)
	}

	if len(stack.Cards) < 3 {
		panic("Not enough cards")
	}
	return stack
}
