package voting

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
	// room *VotingRoom
	id string
}

// func (c *Client) Read() {
// 	defer func() {
// 		c.room.mu.Lock()
// 		delete(c.room.clients, c.conn)
// 		c.room.mu.Unlock()

// 		close(c.send)
// 	}()

// 	for {
// 		_, msg, err := c.conn.ReadMessage()
// 		if err != nil {
// 			return
// 		}

// 		c.room.mu.Lock()
// 		// for client := range c.room.clients {
// 		// 	if client != c.conn {
// 		// 		select {
// 		// 		case client.send <- msg:
// 		// 		default:
// 		// 			close(client.send)
// 		// 			delete(c.room.clients, client)
// 		// 		}
// 		// 	}
// 		// }
// 		c.room.mu.Unlock()
// 	}
// }

func (c *Client) Write() {
	defer func() {
		c.conn.Close()
	}()

	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}

func (c *Client) JoinRoom(r *VotingRoom) {
	// c.room = r
	r.mu.Lock()
	r.clients[c.id] = c
	r.mu.Unlock()
}

func NewClient(conn *websocket.Conn) (*Client, string) {
	id := uuid.NewString()
	c := &Client{
		conn: conn,
		send: make(chan []byte, 2048),
		// room: room,
		id: id,
	}
	return c, id
}
