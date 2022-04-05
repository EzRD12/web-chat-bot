package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ezrod12/chat/messager"
	"github.com/ezrod12/chat/models"
	"github.com/ezrod12/chat/services"
	"github.com/gorilla/websocket"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
	StockPattern   = `/stock=(?P<Stock>.*)`
	JoinPattern    = `/join=(?P<Join>.*)`
	QuitPattern    = `/quit=(?P<Quit>.*)`
)

var (
	newline  = []byte{'\n'}
	space    = []byte{' '}
	mapMutex = sync.RWMutex{}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	messageCollection *mongo.Collection

	// Nickname.
	nick string

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Options to be send to the server
	options chan<- Option
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, result, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		result = bytes.TrimSpace(bytes.Replace(result, newline, space, -1))

		messageRequest := models.Message{}

		json.Unmarshal(result, &messageRequest)

		message := result

		paramsMap := GetParams(JoinPattern, string(message))
		joinKey := "Join"
		mapMutex.RLock()
		if v, ok := paramsMap[joinKey]; ok {
			mapMutex.RUnlock()
			c.hub.unregister <- c
			argument := v
			option := &Option{ID: OPT_JOIN, Client: c, Argument: argument}
			c.options <- *option
			delete(paramsMap, joinKey)
			continue
		}

		paramsMap = GetParams(QuitPattern, string(message))
		quitKey := "Quit"
		mapMutex.RLock()
		if v, ok := paramsMap[quitKey]; ok {
			mapMutex.RUnlock()
			c.hub.unregister <- c
			option := &Option{ID: OPT_QUIT, Client: c, Argument: v}
			c.options <- *option
			delete(paramsMap, quitKey)
			continue
		}

		c.hub.broadcast <- result

		hours, minutes, _ := time.Now().Clock()
		msg := fmt.Sprintf("%d:%02d - %s", hours, minutes, message)
		AddCurrentMessages(RoomsMessages, c.hub.name, msg)
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

			w.Write(message)

			messageRequest := models.Message{}

			json.Unmarshal(message, &messageRequest)

			fmt.Println(messageRequest)

			if messageRequest.Username != "system" && messageRequest.Username != "robot" {
				go services.AddMessage(messageRequest, c.messageCollection, context.TODO())
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

			paramsMap := GetParams(StockPattern, string(message))
			stockKey := "Stock"
			mapMutex.RLock()
			if v, ok := paramsMap[stockKey]; ok {
				mapMutex.RUnlock()
				message := models.StockMessage{HubName: c.hub.name, ClientRemoteAddress: c.conn.RemoteAddr().String(), Message: v}
				messager.SendMessage(&message)
				delete(paramsMap, stockKey)
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
