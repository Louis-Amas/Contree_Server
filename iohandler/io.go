package iohandler

import (
	"contree/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

var server *socketio.Server

// Context socket io context
type Context struct {
	Connected     bool         `json:"connected"`
	CurrentUser   *models.User `json:"user"`
	CurrentLobby  *Lobby
	CurrentPlayer *Player
}

func checkIfUserIsIdentified(s socketio.Conn) error {
	context := s.Context().(*Context)
	if !context.Connected {
		return errors.New("Not connected")
	}
	return nil
}

func errorToJSONString(err error) string {

	return "{\"error\": \"" + err.Error() + "\"}"
}

func errorBadRequest() string {
	return errorToJSONString(errors.New("Bad request"))
}

func connectHandler(s socketio.Conn) error {
	context := Context{Connected: false, CurrentUser: nil}
	s.SetContext(&context)
	fmt.Println("connected:", s.ID())
	return nil
}

func identityHandler(s socketio.Conn, msg string) string {
	var usr models.User
	json.Unmarshal([]byte(msg), &usr)

	context := s.Context().(*Context)
	if usr.VerifyUser() {
		context.Connected = true
		usr.Socket = s
		context.CurrentUser = &usr
	}
	resp, _ := json.Marshal(context)
	return string(resp)
}

func disconnectHandler(s socketio.Conn, reason string) {
	if s.Context() == nil {
		return
	}
	context := s.Context().(*Context)

	// If user in lobby when disconnect then remove it
	if context.CurrentLobby != nil {
		context.CurrentLobby.removeUser(context.CurrentUser)
	}
	if context.CurrentPlayer != nil {
		context.CurrentPlayer.leaveGame()
	}
}
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

		next.ServeHTTP(w, r)
	})
}

// InitSocketIo init socket io server
func InitSocketIo(router *gin.Engine) *socketio.Server {
	var err error
	server, err = socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	go server.Serve()

	router.GET("/socket.io/*any", gin.WrapH(corsMiddleware(server)))
	router.POST("/socket.io/*any", gin.WrapH(corsMiddleware(server)))
	applyEndpoints()
	initLobbies()
	return server
}

// ApplyEndpoints set all handers of socketIo server
func applyEndpoints() {
	server.OnConnect("/", connectHandler)

	server.OnDisconnect("/", disconnectHandler)

	// Identify
	server.OnEvent(DefaultNamespace, OnIdentity, identityHandler)
	// Lobby
	server.OnEvent(DefaultNamespace, OnJoin, joinLobbyHandler)
	server.OnEvent(DefaultNamespace, OnChoosePosition, choosePositionHandler)
	server.OnEvent(DefaultNamespace, OnLeaveLobby, leaveLobbyHandler)

	// Game
	server.OnEvent(DefaultNamespace, OnLeaveGame, leaveGameHandler)

}

func broadcastToRoom(roomName string, event string, msg string) {
	server.BroadcastToRoom(DefaultNamespace, roomName, event, msg)
}
