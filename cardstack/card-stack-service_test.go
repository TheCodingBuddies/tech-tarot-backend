package cardstack

import (
	"sync"
	"testing"
)

func TestSuccessfulInit(t *testing.T) {
	resetService()
	err := Init("../assets/cardstackTest.json")
	if err != nil && !isInit {
		t.Errorf("Cardstack Service was not initialized")
	}
}

func TestServiceIsNotInitialized(t *testing.T) {
	resetService()
	panicOnNoInit := func() {
		GetService()
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	panicOnNoInit()
}

func TestErrorOnMissingStackFile(t *testing.T) {
	resetService()
	err := Init("../notAvailable")
	if err == nil && isInit {
		t.Errorf("Cardstack Service was wrongly initialized")
	}
}

func TestDraw3DifferentCards(t *testing.T) {
	resetService()
	_ = Init("../assets/cardstackTest.json")
	drawnCards := GetService().Draw3()
	if len(drawnCards) != 3 {
		t.Errorf("len(stack.cards) = %d; want 3", len(drawnCards))
	}
	if drawnCards[0].Name == drawnCards[1].Name ||
		drawnCards[0].Name == drawnCards[2].Name ||
		drawnCards[1].Name == drawnCards[2].Name {
		t.Errorf("at least 2 cards are equal!")
	}
}

func resetService() {
	isInit = false
	once = sync.Once{}
}
