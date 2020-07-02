package iohandler

import (
	"contree/game"
	"contree/models"
	"encoding/json"

	"github.com/google/uuid"
	socketio "github.com/googollee/go-socket.io"
)

// ConcreteForwarder struct
type ConcreteForwarder struct {
	gameID string
}

func (fw *ConcreteForwarder) sendToRoom(event string, msg []byte) {
	broadcastToRoom(fw.gameID, event, string(msg))
}

// ForwardStartGame send start game message
func (fw *ConcreteForwarder) ForwardStartGame(game *game.Game) {
	str, _ := json.Marshal(game)
	fw.sendToRoom(EmitStartGame, str)
}

type startCardsMsg struct {
	BestCall *game.Call `json:"bestCall"`
	Team     string     `json:"team"`
}

// ForwardStartCards send start cards message to all clients
func (fw *ConcreteForwarder) ForwardStartCards(team string, bestCall *game.Call) {
	msgStruct := startCardsMsg{BestCall: bestCall, Team: team}
	msg, _ := json.Marshal(&msgStruct)
	fw.sendToRoom(EmitStartCard, msg)
}

type callMsg struct {
	Call     *game.Call `json:"call"`
	Position string     `json:"position"`
}

// ForwardCall forward call to game players
func (fw *ConcreteForwarder) ForwardCall(position string, call *game.Call) {
	msgStruct := callMsg{Call: call, Position: position}
	msg, _ := json.Marshal(&msgStruct)
	fw.sendToRoom(EmitCall, msg)
}

type cardMsg struct {
	Card     *game.Card `json:"card"`
	Position string     `json:"position"`
}

// ForwardCard forward card to game players
func (fw *ConcreteForwarder) ForwardCard(position string, card *game.Card) {
	msgStruct := cardMsg{Position: position, Card: card}
	msg, _ := json.Marshal(&msgStruct)
	fw.sendToRoom(EmitCard, msg)
}

type beloteMsg struct {
	IsReBelote bool   `json:"isReBelote"`
	Position   string `json:"position"`
}

// ForwardBelote forward belote to game
func (fw *ConcreteForwarder) ForwardBelote(position string, isReBelote bool) {
	msgStruct := beloteMsg{Position: position, IsReBelote: isReBelote}
	msg, _ := json.Marshal(&msgStruct)
	fw.sendToRoom(EmitBelote, msg)
}

type positionMsg struct {
	Position string `json:"position"`
}

// ForwardBeginnerPos forward beginner position to player
func (fw *ConcreteForwarder) ForwardBeginnerPos(position string) {
	msgStruct := positionMsg{Position: position}
	msg, _ := json.Marshal(&msgStruct)
	fw.sendToRoom(EmitBeginnerPos, msg)
}

type scoreMsg struct {
	NSScore int `json:"nsScore"`
	WEScore int `json:"weScore"`
}

// ForwardPliScore forward teams scores to players
func (fw *ConcreteForwarder) ForwardPliScore(NSScore, WEScore int) {
	msgStruct := scoreMsg{NSScore: NSScore, WEScore: WEScore}
	msg, _ := json.Marshal(&msgStruct)
	fw.sendToRoom(EmitPliScores, msg)
}

// ForwardWinningPliPosition forward winning position of pli
func (fw *ConcreteForwarder) ForwardWinningPliPosition(position string) {
	msgStruct := positionMsg{Position: position}
	msg, _ := json.Marshal(&msgStruct)
	fw.sendToRoom(EmitWinningPliPosition, msg)
}

// ForwardScore foward global score
func (fw *ConcreteForwarder) ForwardScore(NSScore, WEScore int) {
	msgStruct := scoreMsg{NSScore: NSScore, WEScore: WEScore}
	msg, _ := json.Marshal(&msgStruct)
	fw.sendToRoom(EmitScores, msg)
}

// ForwardEndOfGame end of game
func (fw *ConcreteForwarder) ForwardEndOfGame() {
	fw.sendToRoom(EmitEndOfGame, []byte(""))
}

func leaveGameHandler(s socketio.Conn, msg string) string {
	context := s.Context().(*Context)
	if context.CurrentPlayer != nil {
		context.CurrentPlayer.leaveGame()
		context.CurrentPlayer = nil
	}
	return msg
}

// launch a new go routine init game attributes and start a game
func createGame(north *models.User, west *models.User, south *models.User, east *models.User) {
	go func() {
		gameID := uuid.New().String()

		abort := make(chan bool)
		fw := ConcreteForwarder{gameID: gameID}
		n := newPlayer(north, "NORTH", abort)
		n.joinGame(gameID)
		w := newPlayer(west, "WEST", abort)
		w.joinGame(gameID)

		s := newPlayer(south, "SOUTH", abort)
		s.joinGame(gameID)

		e := newPlayer(east, "EAST", abort)
		e.joinGame(gameID)

		game := game.NewGame(&fw, n, w, s, e)
		game.StartGame()
	}()
}
