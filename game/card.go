package game

// Families cards family
var Families = make([]string, 4)
var normalValue = make(map[string]int)
var assetValue = make(map[string]int)

func initCard() {
	Families[0] = "clovers"
	Families[1] = "hearts"
	Families[2] = "diamonds"
	Families[3] = "spades"

	normalValue["A"] = 11
	normalValue["10"] = 10
	normalValue["K"] = 4
	normalValue["Q"] = 3
	normalValue["J"] = 2
	normalValue["9"] = 0
	normalValue["8"] = 0
	normalValue["7"] = 0

	assetValue["J"] = 20
	assetValue["9"] = 14
	assetValue["A"] = 11
	assetValue["10"] = 10
	assetValue["K"] = 4
	assetValue["Q"] = 3
	assetValue["8"] = 0
	assetValue["7"] = 0

}

// Card struct
type Card struct {
	Value  string `json:"value"`
	Family string `json:"family"`
}

func (c *Card) computeValue(asset string) int {
	if c.Family == asset {
		value, _ := assetValue[c.Value]
		return value
	}
	value, _ := normalValue[c.Value]
	return value
}

func (c *Card) EqualsTo(card *Card) bool {
	return c.Family == card.Family && c.Value == card.Value
}

func (c *Card) compareTo(card2 *Card, asset string) int {
	res := c.computeValue(asset) - card2.computeValue(asset)
	return res
}

// computeBestCardOfPli return index and best card of pli
func computeBestCardOfPli(pli []*Card, asset string) (int, *Card) {

	bestIndex := 0
	bestCard := pli[0]

	for index, card := range pli[1:] {
		if bestCard.Family == card.Family {
			if bestCard.compareTo(card, asset) < 0 {
				bestCard = card
				bestIndex = index + 1
			}
		} else if card.Family == asset {
			bestCard = card
			bestIndex = index + 1
		}
	}
	return bestIndex, bestCard
}

func getAllFamilyCardFromCard(cards []*Card, family string) []*Card {
	newCards := make([]*Card, 0)
	for _, card := range cards {
		if card.Family == family {
			newCards = append(newCards, card)
		}
	}
	return newCards
}

func getBetterAsset(cardsFromFamilies []*Card, bestCard *Card, asset string) []*Card {
	availableCards := make([]*Card, 0)
	for _, card := range cardsFromFamilies {
		if bestCard.compareTo(card, asset) < 0 {
			availableCards = append(availableCards, card)
		}
	}
	return availableCards
}

// ComputeAvailableCards return authorized cards
func ComputeAvailableCards(cards []*Card, pli []*Card, asset string) []*Card {
	if len(pli) == 0 {
		return cards
	}

	askedCard := pli[0]

	// always return normal cards before asset
	if askedCard.Family != asset {
		cardsFromFamilies := getAllFamilyCardFromCard(cards, askedCard.Family)
		if len(cardsFromFamilies) != 0 {
			return cardsFromFamilies
		}
	}

	bestCardIndex, bestCard := computeBestCardOfPli(pli, asset)

	cardsFromFamilies := getAllFamilyCardFromCard(cards, bestCard.Family)

	// already cut or asked card is asset
	if bestCard.Family == asset {
		betterCards := getBetterAsset(cardsFromFamilies, bestCard, asset)

		// if team mate win no need to cut needed
		if len(pli) == 2 && bestCardIndex == 0 && askedCard.Family != asset {
			return cards
		}
		// if team mate win no need to cut needed
		if len(pli) == 3 && bestCardIndex == 1 && askedCard.Family != asset {
			return cards
		}
		// if no other better asset
		if len(cardsFromFamilies) != 0 && len(betterCards) == 0 {
			// if asked card is asset player need to play asset
			if askedCard.Family == asset {
				return cardsFromFamilies
			}
			// can play other cards
			return cards
		}
		// if better assets return
		if len(betterCards) != 0 {
			return betterCards
		}
	}

	// if no card from asked family / best card
	// find asset
	assetCards := getAllFamilyCardFromCard(cards, asset)

	// if team mate win no need to cut needed
	if len(pli) == 2 && bestCardIndex == 0 {
		return cards
	}
	// if team mate win no need to cut needed
	if len(pli) == 3 && bestCardIndex == 1 {
		return cards
	}
	if len(assetCards) != 0 {
		return assetCards
	}
	return cards

}

func computePointOfPli(pli []*Card, asset string) int {
	points := 0
	for _, card := range pli {
		points += card.computeValue(asset)
	}
	return points
}

func computeHasBelote(cards []*Card, asset string) bool {
	count := 0
	for _, c := range cards {
		if c.Family == asset && (c.Value == "K" || c.Value == "Q") {
			count++
		}
	}
	return count == 2
}
