package game

import (
	"math/rand"
)

// Deck struct
type Deck struct {
	cards []*Card
}

func newDeck() *Deck {
	deck := Deck{}
	deck.cards = make([]*Card, 32)

	for indexFamily, family := range Families {
		currentCardIndex := 0
		for key := range normalValue {
			deck.cards[indexFamily*8+currentCardIndex] = &Card{Family: family, Value: key}
			currentCardIndex++
		}
	}
	deck.shuffle()

	return &deck
}

func (deck *Deck) shuffle() {
	dest := make([]*Card, len(deck.cards))
	perm := rand.Perm(len(deck.cards))
	for currentIndex, newIndex := range perm {
		dest[newIndex] = deck.cards[currentIndex]
	}
	deck.cards = dest
}

func (deck *Deck) get4SlicesOf8Cards() [][]*Card {
	slices := make([][]*Card, 4)
	slices[0] = deck.cards[:8]
	slices[1] = deck.cards[8:16]
	slices[2] = deck.cards[16:24]
	slices[3] = deck.cards[24:32]

	return slices
}
