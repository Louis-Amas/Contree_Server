package iohandler

import (
	"contree/game"
	"contree/models"
	"encoding/json"
	"errors"
)

// Player struct socketio wrapper
type Player struct {
	User     *models.User `json:"user"`
	Hand     []*game.Card `json:"hand"`
	Position string       `json:"position"`
	abort    chan bool
}

func newPlayer(u *models.User, pos string, abort chan bool) *Player {
	return &Player{User: u, Position: pos, abort: abort}
}

// GetPosition retur player position
func (player *Player) GetPosition() string {
	return player.Position
}

func (player *Player) joinGame(gameID string) {
	player.User.Socket.Join(gameID)
	c := player.User.Socket.Context().(*Context)
	c.CurrentPlayer = player
}

func (player *Player) leaveGame() {
	player.abort <- true
}

// SendHand send hand to client
func (player *Player) SendHand(cards []*game.Card) {

	player.Hand = cards
	str, _ := json.Marshal(cards)

	player.User.Socket.Emit(EmitHand, string(str))
}

// GetHand return Plyer hand
func (player *Player) GetHand() []*game.Card {
	return player.Hand
}

// AskCall ask call to client (Launch a new go routine and wait until
// new go routine send back error or call
func (player *Player) AskCall() (*game.Call, error) {
	call := make(chan *game.Call)
	errs := make(chan error)

	player.User.Socket.Emit(EmitAskCall, func(data string) {

		var receivedCall game.Call
		err := json.Unmarshal([]byte(data), &receivedCall)
		if err != nil {
			errs <- err
			return
		}
		call <- &receivedCall
	})
	select {
	case c := <-call:
		return c, nil
	case err := <-errs:
		return nil, err
	case <-player.abort:
		return nil, errors.New("Game end")
	}
}

// AskCard ask card to client
func (player *Player) AskCard(availableCards []*game.Card) (*game.Card, error) {
	card := make(chan *game.Card)
	errs := make(chan error)

	encodedCards, _ := json.Marshal(availableCards)

	player.User.Socket.Emit(EmitAskCard, string(encodedCards), func(data string) {
		var receivedCard game.Card
		err := json.Unmarshal([]byte(data), &receivedCard)
		if err != nil {
			errs <- err
			return
		}
		card <- &receivedCard
	})
	select {
	case c := <-card:
		newHand := make([]*game.Card, 0)
		for _, ca := range player.Hand {
			if ca.EqualsTo(c) {
				continue
			}
			newHand = append(newHand, ca)
		}
		player.Hand = newHand
		return c, nil
	case err := <-errs:
		return nil, err
	case <-player.abort:
		return nil, errors.New("Game end")
	}

}
