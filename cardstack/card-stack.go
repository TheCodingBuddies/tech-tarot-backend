package cardstack

import (
	"math/rand"
	"strconv"
)

type Card struct {
	Name string `json:"name"`
}

type Stack struct {
	cards []Card
}

func (stack *Stack) Draw3() []Card {
	stack.shuffle()
	return stack.cards[:3]
}

func (stack *Stack) shuffle() {
	rand.Shuffle(len(stack.cards), func(i, j int) {
		stack.cards[i], stack.cards[j] = stack.cards[j], stack.cards[i]
	})
}

func NewStack() *Stack {
	stack := &Stack{}
	stack.cards = make([]Card, 10)
	for i := range stack.cards {
		stack.cards[i].Name = "Card" + strconv.Itoa(i)
	}
	return stack
}
