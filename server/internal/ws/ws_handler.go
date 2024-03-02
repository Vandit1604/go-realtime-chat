package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type JoinRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(g *gin.Context) {
	var req CreateRoomReq

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	g.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(g *gin.Context) {
	conn, err := upgrader.Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	/*
		our path will look something like this
		localhost:8080/ws/joinRoom:roomId?{userId}=1&{username}=user
	*/

	roomId := g.Param("roomId")
	clientId := g.Query("userId")
	username := g.Query("username")

	client := &Client{
		Conn:     conn,
		Message:  make(chan *Message),
		ID:       clientId,
		RoomId:   roomId,
		Username: username,
	}

	msg := &Message{
		Content:  "User %d joined the room",
		RoomId:   roomId,
		Username: username,
	}

	h.hub.Register <- client
	h.hub.Broadcast <- msg

	client.ReadMessage(h.hub)
	go client.WriteMessage()
}

type GetRoomsRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetRooms(ginContext gin.Context) {

}
