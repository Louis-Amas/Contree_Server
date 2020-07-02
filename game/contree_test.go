package game

import (
	"os"
	"testing"
)

func TestCreatePlayerOrder(t *testing.T) {
	correctOrder := []string{
		"EAST",
		"NORTH",
		"WEST",
		"SOUTH",
	}
	order := createPlayerOrder("EAST")
	for i := 0; i < 4; i++ {
		if correctOrder[i] != order[i] {
			t.Errorf("Bad order")
		}
	}
}

func TestGet8SlicesOfDeck(t *testing.T) {
	InitGame()
	deck := newDeck()
	slices := deck.get4SlicesOf8Cards()
	for i := 0; i < 4; i++ {

		if len(slices[i]) != 8 {
			t.Errorf("Bad slices size got %d", len(slices[i]))
		}
	}
}

func TestMain(m *testing.M) {
	InitGame()
	os.Exit(m.Run())
}
