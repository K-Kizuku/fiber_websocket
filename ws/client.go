package ws

import (
	"fiber_websocket/db"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *fasthttp.RequestCtx) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	room *Room

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	id string

	mu sync.Mutex
}

func (c *Client) readPump() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()
	// var v V
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		slices := strings.Split(string(message), " ")
		c.room.mu.Lock()
		for i, block := range slices {
			switch block {
			case "--M":
				log.Printf("--M: %v", slices[1+i])
				c.room.broadcast <- []byte("--M" + slices[i+1])
			case "--O":
				log.Printf("--O: %v", slices[1+i])
				if c.room.Players[slices[i+1]]++; c.room.Players[slices[i+1]]%10 == 0 {
					for k, v := range c.room.Players {
						db.Db.Model(&db.Point{}).Where("id = ?", k).Update("point", strconv.Itoa(v))
					}

				}
			}
		}
		c.room.mu.Unlock()
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	c.conn.WriteMessage(websocket.TextMessage, []byte(c.id))
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			// c.conn.WriteJSON()
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (room *Room) ServeWs(ctx *fiber.Ctx) error {
	err := upgrader.Upgrade(ctx.Context(), func(conn *websocket.Conn) {
		// u, _ := uuid.NewRandom()
		type Payload struct {
			userId uuid.UUID `json:"id"`
		}
		var payload Payload
		err := ctx.BodyParser(&payload)
		if err != nil {
			ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		// payload := ctx.Body()

		client := &Client{room: room, conn: conn, send: make(chan []byte, 512), id: payload.userId.String()}
		client.room.register <- client
		log.Println("tets")
		log.Printf("uuid: %s", payload.userId.String())

		go client.writePump()
		client.readPump()
	})

	if err != nil {
		log.Println(err)
	}
	return nil
}
