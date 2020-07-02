package game

import (
	"math/rand"
	"time"
)

var positionOrder []string

// Player interface
type Player interface {
	GetPosition() string
	SendHand(cards []*Card)
	GetHand() []*Card
	AskCall() (*Call, error)
	AskCard(availableCards []*Card) (*Card, error)
}

// Forwarder interface
type Forwarder interface {
	ForwardStartGame(game *Game)
	ForwardStartCards(team string, bestCall *Call)
	ForwardCall(position string, call *Call)
	ForwardCard(position string, card *Card)
	ForwardBelote(position string, isReBelote bool)
	ForwardBeginnerPos(position string)
	ForwardWinningPliPosition(position string)
	ForwardPliScore(NSScore, WEScore int)
	ForwardScore(NSScore, WEScore int)
	ForwardEndOfGame()
}

const teamNorthSouth = "NS"
const teamWestEast = "WE"

// Game struct
type Game struct {
	North Player
	West  Player
	South Player
	East  Player

	playersOrder []Player

	// currentPos who is currently Playing
	currentPos string

	// who is currently calling
	callPos string

	winningCallPos string
	winningCall    *Call

	belote   bool
	reBelote bool

	contree      bool
	isSurContree bool

	currentHandScoreNS int
	currentHandScoreWE int

	scoreNS int
	scoreWE int

	fw    Forwarder
	pli   []*Card
	calls []*Call
}

// InitGame package
func InitGame() {
	positionOrder = []string{
		"NORTH",
		"WEST",
		"SOUTH",
		"EAST",
	}
	rand.Seed(time.Now().UTC().UnixNano())
	initCard()

}

func createPlayerOrder(firstPlayerPosition string) []string {

	indexFirstPlayer := 0
	for index, el := range positionOrder {
		if el == firstPlayerPosition {
			indexFirstPlayer = index
			break
		}
	}
	order := make([]string, 4)
	currentIndex := indexFirstPlayer
	for i := 0; i < 4; i++ {
		order[i] = positionOrder[currentIndex]
		if currentIndex == 3 {
			currentIndex = 0
		} else {
			currentIndex++
		}
	}
	return order
}

// NewGame init game
func NewGame(fw Forwarder, north Player, west Player, south Player, east Player) *Game {
	game := Game{}
	game.North = north
	game.West = west
	game.South = south
	game.East = east

	game.calls = make([]*Call, 0)

	game.fw = fw

	return &game
}

func getNextPosition(position string) string {
	switch position {
	case "NORTH":
		return "WEST"
	case "WEST":
		return "SOUTH"
	case "SOUTH":
		return "EAST"
	case "EAST":
		return "NORTH"
	}
	return ""
}

func (game *Game) getPlayerFromPosition(position string) Player {
	if position == "NORTH" {
		return game.North
	}
	if position == "WEST" {
		return game.West
	}
	if position == "SOUTH" {
		return game.South
	}
	if position == "EAST" {
		return game.East
	}
	return nil
}

func (game *Game) getAsset() string {
	return game.winningCall.Family
}

func (game *Game) computeCurrentPositionOrder(position string) {
	order := createPlayerOrder(position)
	game.playersOrder = make([]Player, 4)
	for i, position := range order {
		game.playersOrder[i] = game.getPlayerFromPosition(position)
	}
}

func (game *Game) distribute() {
	deck := newDeck()
	hands := deck.get4SlicesOf8Cards()
	for index, player := range game.playersOrder {
		player.SendHand(hands[index])
	}
}

func (game *Game) askCalls() (bool, error) {
	var winningPlayer Player
	for countPass := 0; countPass != 3; {
		for _, player := range game.playersOrder {
			call, err := player.AskCall()
			if err != nil {
				return false, err
			}
			if game.winningCall == nil && !call.IsPass {
				game.winningCall = call
				winningPlayer = player
				countPass = 0
			} else {

				if call.IsPass {
					countPass++
				} else if call.IsContree || call.IsSurContree {

					if call.IsSurContree {
						countPass = 3
						game.isSurContree = true
						game.contree = false
					} else {
						game.contree = true
					}
				} else {
					game.winningCall = call
					winningPlayer = player
					countPass = 0
				}

			}
			game.calls = append(game.calls, call)
			game.fw.ForwardCall(player.GetPosition(), call)
			if countPass == 3 && game.winningCall == nil {
				continue
			}
			// no call found
			if countPass == 4 {
				return false, nil
			}
			if countPass >= 3 {

				if winningPlayer.GetPosition() == "NORTH" || winningPlayer.GetPosition() == "SOUTH" {
					game.winningCallPos = teamNorthSouth
				} else {
					game.winningCallPos = teamWestEast
				}
				return true, nil
			}
		}
	}
	return false, nil

}

func (game *Game) addPointsToPlayerTeam(player Player, points int) {
	var ptr *int
	if player.GetPosition() == "NORTH" || player.GetPosition() == "SOUTH" {
		ptr = &game.currentHandScoreNS
	} else {
		ptr = &game.currentHandScoreWE
	}
	*ptr += points
}

func (game *Game) manageBelote() {
	if game.winningCallPos == teamNorthSouth {
		if computeHasBelote(game.North.GetHand(), game.getAsset()) {
			game.belote = true
		}
		if computeHasBelote(game.South.GetHand(), game.getAsset()) {
			game.belote = true
		}
	} else {
		if computeHasBelote(game.West.GetHand(), game.getAsset()) {
			game.belote = true
		}
		if computeHasBelote(game.East.GetHand(), game.getAsset()) {
			game.belote = true
		}
	}
}

func (game *Game) checkIfPlayerBelongsToCallingTeam(player Player) bool {
	if game.winningCallPos == teamNorthSouth {
		return player.GetPosition() == "NORTH" || player.GetPosition() == "SOUTH"
	}
	return player.GetPosition() == "WEST" || player.GetPosition() == "EAST"

}

func (game *Game) checkIfPlayerPlayBelote(card *Card) bool {
	return card.Family == game.getAsset() && (card.Value == "K" || card.Value == "Q")
}
func (game *Game) manageBelotePlay(player Player, card *Card) {
	if (game.belote || game.reBelote) && game.checkIfPlayerBelongsToCallingTeam(player) {
		if game.checkIfPlayerPlayBelote(card) {
			// first send and add the point
			if game.belote {
				game.addPointsToPlayerTeam(player, 20)
				game.belote = false
				game.reBelote = true
				game.fw.ForwardBelote(player.GetPosition(), false)
				game.sendPliScores()
			} else if game.reBelote {
				game.fw.ForwardBelote(player.GetPosition(), true)
				game.reBelote = false
			}
		}
	}
}

func (game *Game) askCards() error {
	asset := game.getAsset()
	game.manageBelote()
	for i := 0; i < 8; i++ {
		pli := make([]*Card, 0)
		game.fw.ForwardBeginnerPos(game.playersOrder[0].GetPosition())

		for _, player := range game.playersOrder {

			availableCards := ComputeAvailableCards(player.GetHand(), pli, asset)
			card, err := player.AskCard(availableCards)
			if err != nil {
				return err
			}
			game.manageBelotePlay(player, card)
			pli = append(pli, card)

			game.fw.ForwardCard(player.GetPosition(), card)
		}
		index, _ := computeBestCardOfPli(pli, asset)
		winningPlayer := game.playersOrder[index]
		game.currentPos = winningPlayer.GetPosition()
		game.fw.ForwardWinningPliPosition(winningPlayer.GetPosition())

		game.addPointsToPlayerTeam(winningPlayer, computePointOfPli(pli, asset))

		// Dix de der
		if i == 7 {
			game.addPointsToPlayerTeam(winningPlayer, 10)
		}

		game.sendPliScores()

		game.computeCurrentPositionOrder(game.currentPos)
	}
	return nil
}

func (game *Game) sendPliScores() {
	game.fw.ForwardPliScore(game.currentHandScoreNS, game.currentHandScoreWE)
}

func (game *Game) sendGlobalScores() {
	game.computeScores()
	game.fw.ForwardScore(game.scoreNS, game.scoreWE)
}

func (game *Game) computeScores() {
	coef := 1

	if game.contree {
		coef = 2
	} else if game.isSurContree {
		coef = 4
	}

	var callingTeamCurrentHandScore *int
	var callingTeamScore *int
	var notCallingTeamScore *int
	if game.winningCallPos == teamNorthSouth {
		callingTeamCurrentHandScore = &game.currentHandScoreNS
		callingTeamScore = &game.scoreNS
		notCallingTeamScore = &game.scoreWE
	} else {
		callingTeamCurrentHandScore = &game.currentHandScoreWE
		callingTeamScore = &game.scoreWE
		notCallingTeamScore = &game.scoreNS
	}
	bestCall := game.winningCall

	callMaxPoint := bestCall.Value
	neededPointToWin := bestCall.Value
	opponentsMaxPoints := 160
	if bestCall.Value == 190 {
		neededPointToWin = 162
		callMaxPoint = 250
		opponentsMaxPoints = 500
	} else if bestCall.Value == 200 {
		neededPointToWin = 182
		callMaxPoint = 270
		opponentsMaxPoints = 520
	}

	if neededPointToWin <= *callingTeamCurrentHandScore {
		*callingTeamScore += callMaxPoint * coef
	} else {
		*notCallingTeamScore += opponentsMaxPoints * coef
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func (game *Game) reset() {

	game.currentHandScoreNS = 0
	game.currentHandScoreWE = 0
	game.contree = false
	game.isSurContree = false
	game.winningCall = nil
	game.winningCallPos = ""
	game.contree = false
	game.isSurContree = false
	game.belote = false
	game.reBelote = false
}

// StartGame start game in current go routine
func (game *Game) StartGame() {
	game.fw.ForwardStartGame(game)
	game.callPos = positionOrder[randInt(0, 3)]
	game.currentPos = game.callPos

	for game.scoreNS < 2000 && game.scoreWE < 2000 {
		game.computeCurrentPositionOrder(game.callPos)
		game.distribute()

		hasFoundCall, err := game.askCalls()
		if err != nil {
			break
		}
		// if found a call
		if !hasFoundCall {
			continue
		}
		game.fw.ForwardStartCards(game.winningCallPos, game.winningCall)

		if game.askCards() != nil {
			break
		}
		game.sendGlobalScores()

		game.reset()

		game.callPos = getNextPosition(game.callPos)
	}

	game.fw.ForwardEndOfGame()
	return
}
