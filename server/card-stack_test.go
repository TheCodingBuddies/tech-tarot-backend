package server

import "testing"

func TestDraw3Cards(t *testing.T) {
	stack := NewStack()
	drawnCards := stack.Draw3()
	if len(drawnCards) != 3 {
		t.Errorf("len(stack.cards) = %d; want 3", len(drawnCards))
	}
}

func TestDraw3DifferentCards(t *testing.T) {
	drawnCards := NewStack().Draw3()
	if drawnCards[0].Name == drawnCards[1].Name ||
		drawnCards[1].Name == drawnCards[2].Name ||
		drawnCards[0].Name == drawnCards[2].Name {
		t.Errorf("at least 2 cards are equal!")
	}
}
