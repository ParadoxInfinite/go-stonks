package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 1 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the exchange.
type Client struct {
	exch *Exchange

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	id string
}

// readPump pumps messages from the websocket connection to the exchange.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.exch.unregister <- c
		c.conn.Close()
	}()
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
		log.Printf("Incoming req: %s", message)
	}
}

// writePump pumps messages from the exchange to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	tick := time.NewTicker(writeWait)
	defer func() {
		ticker.Stop()
		tick.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case _, ok := <-c.send:
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
			message, _ := json.Marshal("Currently we only stream stocks data")
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-tick.C:
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Error: %s, while sending message to %s", err, c.conn.RemoteAddr())
				return
			}
			var stocks []stock
			collection := CNX.Database("go-stonks").Collection("stonks")
			collection.Find(CTX, bson.D{})
			cur, err := collection.Find(CTX, bson.D{})
			if err != nil {
				log.Fatal(err)
			}
			defer cur.Close(CTX)
			for cur.Next(CTX) {
				var result stock
				err := cur.Decode(&result)
				if err != nil {
					log.Fatal(err)
				}
				stocks = append(stocks, result)
			}
			var updatedStocks []stock
			for _, element := range stocks {
				var change float64
				if rand.Float64() > 0.5 {
					change = rand.Float64()
				} else {
					change = -1 * rand.Float64()
				}
				currentStock := stock{
					ID:       element.ID,
					Name:     element.Name,
					Symbol:   element.Symbol,
					Exchange: element.Exchange,
					Price:    element.Price + change,
				}
				log.Println(element)
				insertResult, err := collection.UpdateOne(CTX, bson.M{"id": element.ID}, bson.M{"$set": bson.M{"price": currentStock.Price}})
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("Successfully updated %s stock! Result: %d", currentStock.Name, insertResult.ModifiedCount)
				updatedStocks = append(updatedStocks, currentStock)
			}
			log.Println(updatedStocks)
			stonks, _ := json.Marshal(updatedStocks)
			w.Write(stonks)
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(exch *Exchange, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// TODO: add some sort of authentication and display stocks only for that user.
	id := r.Header.Get("id")
	// if id == "" {
	// 	log.Printf("Client did not send id. IP: %s", conn.RemoteAddr())
	// 	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "You did not send an ID."))
	// }
	client := &Client{exch: exch, conn: conn, send: make(chan []byte, 256), id: id}
	client.exch.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
