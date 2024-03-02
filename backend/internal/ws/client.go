package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	RoomId   string `json:"roomid"`
	Username string `json:"username"`
}

type Message struct {
	Content  string `json:"content"`
	RoomId   string `json:"roomid"`
	Username string `json:"sentby"`
}

func (c *Client) WriteMessage() {
	// After writing message close the websocket connection
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}
		c.Conn.WriteJSON(message)
	}
}

func (c *Client) ReadMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()

	}()

	for {
		_, message, err := c.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &Message{
			Content:  string(message),
			RoomId:   c.RoomId,
			Username: c.Username,
		}

		hub.Broadcast <- msg

	}
}
