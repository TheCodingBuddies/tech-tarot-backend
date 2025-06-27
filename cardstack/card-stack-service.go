package cardstack

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
)

type StackService struct {
	stack *Stack
}

var (
	instance  *StackService
	once      sync.Once
	initError error
	isInit    bool
)

func GetService() *StackService {
	if !isInit {
		panic("Cardstack Service was not initialized")
	}
	return instance
}

func Init(path string) error {
	once.Do(func() {
		stack := &Stack{}
		loadedStackFile, err := os.Open(path)
		if err != nil {
			initError = fmt.Errorf("failed to load stack file: %w", err)
			return
		}

		defer loadedStackFile.Close()

		err = json.NewDecoder(loadedStackFile).Decode(stack)
		if err != nil {
			initError = fmt.Errorf("failed to decode stack: %w", err)
			return
		}

		if len(stack.Cards) < 3 {
			initError = fmt.Errorf("not enough cards: found %d, require at least 3", len(stack.Cards))
			return
		}

		instance = &StackService{stack: stack}
		isInit = true
	})
	return initError
}

func (service *StackService) GetImageURL(id int) (string, error) {
	card, exists := service.stack.Cards[id]
	if !exists {
		return "", fmt.Errorf("card with id %d not found", id)
	}
	return card.ImageSrc, nil
}

func (service *StackService) Draw3() []CardDto {
	shuffledCards := service.shuffle()
	return shuffledCards[:3]
}

func (service *StackService) shuffle() []CardDto {
	var cards []CardDto
	for i, card := range service.stack.Cards {
		cards = append(cards, CardDto{
			Id:   i,
			Name: card.Name,
		})
	}
	// very important information: fisher-yates algo
	for i := len(cards) - 1; i > 0; i-- {
		randPos := rand.Intn(i + 1)
		cards[i], cards[randPos] = cards[randPos], cards[i]
	}

	return cards
}
