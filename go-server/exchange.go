package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// Stock structure
type stock struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Price    float64 `json:"price"`
	Exchange string  `json:"exchange"` // TODO: OPTIONAL: add multiple exchanges
}

// TODO: This is v0.0.1 sample data. Will have to migrate this to DB in future commits.
var stocks = []stock{
	{ID: "stonkschange_00000001", Name: "C", Symbol: "C", Price: 3.14, Exchange: "STNKS"},
	{ID: "stonkschange_00000010", Name: "C++", Symbol: "CPP", Price: 420.00, Exchange: "STNKS"},
	{ID: "stonkschange_00000011", Name: "JavaScript", Symbol: "JS", Price: 69.69, Exchange: "STNKS"},
	{ID: "stonkschange_00000100", Name: "TypeScript", Symbol: "TS", Price: 9001.00, Exchange: "STNKS"},
}

type Exchange struct {
	// Registered clients.
	clients map[string]bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newExchange() *Exchange {
	return &Exchange{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]bool),
	}
}

func (e *Exchange) run() {
	for {
		select {
		case client := <-e.register:
			if e.clients[client.id] {
				client.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(1008, "Invalid ID provided."))
			}
			log.Printf("Registering client: %s with IP: %s", client.id, client.conn.RemoteAddr())
			e.clients[client.id] = true
		case client := <-e.unregister:
			if _, ok := e.clients[client.id]; ok {
				log.Printf("Unregistering client: %s with IP: %s", client.id, client.conn.RemoteAddr())
				delete(e.clients, client.id)
				close(client.send)
			}
		}
	}
}
