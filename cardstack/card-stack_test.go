package cardstack

import "testing"

func TestStackFileIsMissing(t *testing.T) {
	panicOnStackLoading := func() {
		NewStack("../is/missing.json")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	panicOnStackLoading()
}

func TestDraw3Cards(t *testing.T) {
	stack := NewStack("../assets/cardstackTest.json")
	drawnCards := stack.Draw3()
	if len(drawnCards) != 3 {
		t.Errorf("len(stack.cards) = %d; want 3", len(drawnCards))
	}
}

func TestDraw3DifferentCards(t *testing.T) {
	drawnCards := NewStack("../assets/cardstackTest.json").Draw3()
	if drawnCards[0].Name == drawnCards[1].Name ||
		drawnCards[1].Name == drawnCards[2].Name ||
		drawnCards[0].Name == drawnCards[2].Name {
		t.Errorf("at least 2 cards are equal!")
	}
}
