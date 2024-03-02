package ws

type Hub struct {
	Rooms      map[string]*Room `json:"rooms"`
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			// check for the Room id then check if the client id exists if yes add the client to the rooms map
			if _, ok := h.Rooms[client.RoomId]; ok {
				room := h.Rooms[client.RoomId]
				if _, ok := room.Clients[client.ID]; ok {
					room.Clients[client.ID] = client
				}
			}
		case client := <-h.Unregister:
			// check for the Room id then check if the client id exists in that room delete the client ID from that room
			if _, ok := h.Rooms[client.RoomId]; ok {
				if _, ok := h.Rooms[client.RoomId].Clients[client.ID]; ok {
					// if more than 1 user in the room notify them
					if len(h.Rooms[client.RoomId].Clients) > 0 {
						h.Broadcast <- &Message{
							Content:  "User left the Room",
							RoomId:   client.RoomId,
							Username: client.Username,
						}
					}
					// send message that a client left the room
					delete(h.Rooms[client.RoomId].Clients, client.ID)
					close(client.Message)
				}
			}

		case message := <-h.Broadcast:
			if _, ok := h.Rooms[message.RoomId]; ok {
				for _, client := range h.Rooms[message.RoomId].Clients {
					client.Message <- message
				}
			}
		}
	}
}
