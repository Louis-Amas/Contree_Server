package game

import "testing"

func TestComputeBestCardOfPli(t *testing.T) {
	pli := []*Card{
		&Card{Value: "J", Family: "hearts"},
		&Card{Value: "Q", Family: "hearts"},
		&Card{Value: "7", Family: "spades"},
		&Card{Value: "10", Family: "hearts"},
	}

	if i, bestCard := computeBestCardOfPli(pli, "spades"); bestCard.Family != "spades" && i != 2 {
		t.Error("Bad winning card")
	}
}
func TestComputeBestCardOfPli_2(t *testing.T) {
	pli := []*Card{
		&Card{Value: "J", Family: "hearts"},
		&Card{Value: "Q", Family: "hearts"},
		&Card{Value: "10", Family: "hearts"},
		&Card{Value: "7", Family: "hearts"},
	}

	if i, bestCard := computeBestCardOfPli(pli, "spades"); bestCard.Value != "10" && i != 2 {
		t.Error("Bad winning card")
	}
}
func TestComputeBestCardOfPli_3(t *testing.T) {
	pli := []*Card{
		&Card{Value: "J", Family: "hearts"},
		&Card{Value: "Q", Family: "hearts"},
		&Card{Value: "10", Family: "hearts"},
		&Card{Value: "7", Family: "hearts"},
	}

	if i, bestCard := computeBestCardOfPli(pli, "hearts"); bestCard.Value != "J" && i != 0 {
		t.Error("Bad winning card")
	}
}
func TestComputeAvailableCards(t *testing.T) {
	pli := []*Card{
		&Card{Value: "J", Family: "hearts"},
		&Card{Value: "Q", Family: "hearts"},
		&Card{Value: "10", Family: "hearts"},
	}
	hand := []*Card{
		&Card{Value: "A", Family: "clovers"},
		&Card{Value: "7", Family: "spades"},
		&Card{Value: "A", Family: "hearts"},
	}

	cards := ComputeAvailableCards(hand, pli, "spades")
	if len(cards) != 1 && cards[0].Value != "A" {
		t.Errorf("Bad available cards")
	}

}
func TestComputeAvailableCards_cut(t *testing.T) {
	pli := []*Card{
		&Card{Value: "J", Family: "hearts"},
		&Card{Value: "Q", Family: "hearts"},
		&Card{Value: "10", Family: "hearts"},
	}
	hand := []*Card{
		&Card{Value: "A", Family: "clovers"},
		&Card{Value: "7", Family: "spades"},
		&Card{Value: "7", Family: "diamonds"},
	}

	cards := ComputeAvailableCards(hand, pli, "spades")
	if len(cards) != 1 && cards[0].Value != "7" && cards[0].Family != "spades" {
		t.Errorf("Bad available cards")
	}

}

func TestComputeAvailableCards_higher_cut(t *testing.T) {
	pli := []*Card{
		&Card{Value: "J", Family: "hearts"},
		&Card{Value: "Q", Family: "hearts"},
		&Card{Value: "10", Family: "spades"},
	}
	hand := []*Card{
		&Card{Value: "A", Family: "clovers"},
		&Card{Value: "J", Family: "spades"},
		&Card{Value: "7", Family: "spades"},
	}

	cards := ComputeAvailableCards(hand, pli, "spades")
	if len(cards) != 1 && cards[0].Value != "J" && cards[0].Family != "spades" {
		t.Errorf("Bad available cards")
	}

}

func TestComputeAvailableCards_cut_but_still_have_normal_card(t *testing.T) {
	pli := []*Card{
		&Card{Value: "J", Family: "hearts"},
		&Card{Value: "7", Family: "spades"},
	}
	hand := []*Card{
		&Card{Value: "A", Family: "hearts"},
		&Card{Value: "J", Family: "spades"},
		&Card{Value: "7", Family: "spades"},
	}

	cards := ComputeAvailableCards(hand, pli, "spades")
	if len(cards) != 1 && cards[0].Value != "A" && cards[0].Family != "hearts" {
		t.Errorf("Bad available cards")
	}

}
func TestComputeAvailableCards_can_cut_but_not_with_higher_value(t *testing.T) {
	pli := []*Card{
		{Value: "K", Family: "hearts"},
		{Value: "9", Family: "clovers"},
	}
	hand := []*Card{
		{Value: "A", Family: "spades"},
		{Value: "7", Family: "spades"},
		{Value: "K", Family: "clovers"},
	}

	cards := ComputeAvailableCards(hand, pli, "clovers")
	if len(cards) != 3 {
		t.Errorf("Bad available cards")
	}

}
func TestComputeAvailableCards_asset_asked_and_has_lower(t *testing.T) {
	pli := []*Card{
		{Value: "A", Family: "hearts"},
		{Value: "Q", Family: "hearts"},
		{Value: "9", Family: "hearts"},
	}
	hand := []*Card{
		{Value: "A", Family: "diamonds"},
		{Value: "8", Family: "hearts"},
	}

	cards := ComputeAvailableCards(hand, pli, "hearts")
	if len(cards) != 1 {
		t.Errorf("Bad available cards")
	}

}

func TestHasBelote_true(t *testing.T) {
	hand := []*Card{
		{Value: "K", Family: "diamonds"},
		{Value: "Q", Family: "diamonds"},
		{Value: "7", Family: "diamonds"},
		{Value: "7", Family: "spades"},
	}

	if !computeHasBelote(hand, "diamonds") {
		t.Errorf("Should have belote")
	}
}

func TestHasBelote_false(t *testing.T) {
	hand := []*Card{
		{Value: "K", Family: "diamonds"},
		{Value: "10", Family: "diamonds"},
		{Value: "7", Family: "diamonds"},
		{Value: "7", Family: "spades"},
	}

	if computeHasBelote(hand, "diamonds") {
		t.Errorf("Should not have belote")
	}
}

func TestHasBelote_team_mate_win_no_need_to_cut(t *testing.T) {
	pli := []*Card{
		{Value: "A", Family: "hearts"},
		{Value: "Q", Family: "diamonds"},
		{Value: "9", Family: "hearts"},
	}
	hand := []*Card{
		{Value: "A", Family: "spades"},
		{Value: "J", Family: "diamonds"},
	}

	cards := ComputeAvailableCards(hand, pli, "diamonds")
	if len(cards) != 2 {
		t.Errorf("Bad available cards")
	}
}
func TestHasBelote_team_mate_win_no_need_to_cut_3(t *testing.T) {
	pli := []*Card{
		{Value: "A", Family: "clovers"},
		{Value: "Q", Family: "spades"},
	}
	hand := []*Card{
		{Value: "8", Family: "spades"},
		{Value: "10", Family: "diamonds"},
		{Value: "J", Family: "diamonds"},
	}

	cards := ComputeAvailableCards(hand, pli, "diamonds")
	if len(cards) != 3 {
		t.Errorf("Bad available cards")
	}
}
func TestHasBelote_team_mate_win_no_need_to_cut_4(t *testing.T) {
	pli := []*Card{
		{Value: "7", Family: "diamonds"},
		{Value: "K", Family: "diamonds"},
	}
	hand := []*Card{
		{Value: "8", Family: "spades"},
		{Value: "10", Family: "diamonds"},
		{Value: "J", Family: "diamonds"},
	}

	cards := ComputeAvailableCards(hand, pli, "diamonds")
	if len(cards) != 2 {
		t.Errorf("Bad available cards")
	}
}
