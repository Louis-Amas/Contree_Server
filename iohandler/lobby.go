package iohandler

import (
	"contree/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	socketio "github.com/googollee/go-socket.io"
)

// North string representing north position
const North = "NORTH"

// West string representing west position
const West = "WEST"

// South string representing south position
const South = "SOUTH"

// East string representing east position
const East = "EAST"

// Lobby struct
type Lobby struct {
	LobbyName        string `json:"lobbyName"`
	positionRequests chan positionRequest
	removeUsers      chan *models.User
	usersCount       int
	North            *models.User `json:"north"`
	West             *models.User `json:"west"`
	South            *models.User `json:"south"`
	East             *models.User `json:"east"`
}

var lobbies map[string]*Lobby
var mutex = &sync.Mutex{}

func initLobbies() {
	lobbies = make(map[string]*Lobby)
}

func joinLobbyHandler(s socketio.Conn, msg string) string {
	err := checkIfUserIsIdentified(s)
	if err != nil {
		return errorToJSONString(err)
	}
	var req Lobby
	// Decode request from json
	err = json.Unmarshal([]byte(msg), &req)
	if err != nil {
		return err.Error()
	}
	mutex.Lock()
	defer mutex.Unlock()
	lobby, ok := lobbies[req.LobbyName]
	// if lobby doesn't exist create it
	if !ok {
		lobby = &req
		lobby.usersCount = 0
		lobbies[lobby.LobbyName] = lobby
		go lobby.handleLobbyChoosePosition()
	}

	s.Join(lobby.LobbyName)

	lobby.usersCount++
	context := s.Context().(*Context)
	context.CurrentLobby = lobby

	return lobby.lobbyToJSONString()
}

func (lobby *Lobby) handleLobbyChoosePosition() {
	lobby.positionRequests = make(chan positionRequest)
	lobby.removeUsers = make(chan *models.User)

	for {
		select {
		case req := <-lobby.positionRequests:

			// Remove user from is current position do nothing if user has no position
			lobby.removeUserFromAllPositions(req.User)
			if !lobby.setPosition(req.Position, req.User) {
				req.Errors <- errors.New("Position already taken")
				continue
			}

			resp := lobby.lobbyToJSONString()

			req.Response <- resp

			lobby.sendToAll(EmitLobby, resp)
			if lobby.isFull() {
				// FIXME: Not good

				fmt.Println("Launching a game")
				lobby.North.Socket.Leave(lobby.LobbyName)
				lobby.West.Socket.Leave(lobby.LobbyName)
				lobby.South.Socket.Leave(lobby.LobbyName)
				lobby.East.Socket.Leave(lobby.LobbyName)
				lobby.North.Socket.Context().(*Context).CurrentLobby = nil
				lobby.West.Socket.Context().(*Context).CurrentLobby = nil
				lobby.South.Socket.Context().(*Context).CurrentLobby = nil
				lobby.East.Socket.Context().(*Context).CurrentLobby = nil
				createGame(lobby.North, lobby.West, lobby.South, lobby.East)

				// delete lobby
				delete(lobbies, lobby.LobbyName)
				// kill the current go routine
				return
			}

		case user := <-lobby.removeUsers:
			mutex.Lock()
			lobby.removeUserFromAllPositions(user)
			lobby.usersCount--
			mutex.Unlock()
			if lobby.isEmpty() {
				delete(lobbies, lobby.LobbyName)
				return
			}
		}
	}
}

type positionRequest struct {
	Position string
	User     *models.User
	Response chan string
	Errors   chan error
}

func (req *positionRequest) fillPositionRequest(user *models.User) {
	req.User = user
	req.Response = make(chan string)
	req.Errors = make(chan error)
}

func choosePositionHandler(s socketio.Conn, msg string) string {
	context := s.Context().(*Context)
	if context.CurrentLobby == nil {
		return errorToJSONString(errors.New("not in a lobby"))
	}
	lobby := context.CurrentLobby
	var req positionRequest
	err := json.Unmarshal([]byte(msg), &req)
	if err != nil {
		return errorBadRequest()
	}
	req.fillPositionRequest(context.CurrentUser)

	lobby.positionRequests <- req

	select {
	case err := <-req.Errors:
		return errorToJSONString(err)
	case resp := <-req.Response:
		return resp
	}
}

func (lobby *Lobby) lobbyToJSONString() string {
	resp, err := json.Marshal(lobby)
	if err != nil {
		log.Fatal(err)
	}
	return string(resp)
}

func (lobby *Lobby) removeUser(user *models.User) {
	lobby.removeUsers <- user
}

func (lobby *Lobby) sendToAll(event string, msg string) {
	broadcastToRoom(lobby.LobbyName, event, msg)
}

func (lobby *Lobby) isFull() bool {
	return lobby.North != nil && lobby.West != nil && lobby.South != nil && lobby.East != nil
}

func (lobby *Lobby) isEmpty() bool {
	return lobby.usersCount == 0
}

func (lobby *Lobby) setPosition(pos string, usr *models.User) bool {

	switch pos {
	case North:
		if lobby.North == nil {
			lobby.North = usr
			return true
		}

	case West:
		if lobby.West == nil {
			lobby.West = usr
			return true
		}

	case South:
		if lobby.South == nil {
			lobby.South = usr
			return true
		}

	case East:
		if lobby.East == nil {
			lobby.East = usr
			return true
		}
	}

	return false
}

func (lobby *Lobby) removeUserFromAllPositions(user *models.User) {
	if lobby.North == user {
		lobby.North = nil
	}
	if lobby.West == user {
		lobby.West = nil
	}
	if lobby.South == user {
		lobby.South = nil
	}
	if lobby.East == user {
		lobby.East = nil
	}

	lobby.sendToAll(EmitLobby, lobby.lobbyToJSONString())
}

func leaveLobbyHandler(s socketio.Conn, msg string) string {
	context := s.Context().(*Context)
	if context.CurrentLobby != nil {
		context.CurrentLobby.removeUser(context.CurrentUser)
		context.CurrentLobby = nil
	}
	return msg
}
