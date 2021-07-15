package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var mutex = sync.Mutex{}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan Speakers
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		var msg Speakers
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		mutex.Lock()
		c.hub.speakers = msg
		mutex.Unlock()
		c.hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	for {
		msg, ok := <-c.send
		if !ok {
			return
		}
		err := c.conn.WriteJSON(msg)
		if err != nil {
			log.Printf("error: %v", err)
			c.conn.Close()
		}
	}
}

func handleConnections(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println("Handle connection")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan Speakers, 256)}
	client.hub.register <- client

	go client.readPump()
	go client.writePump()
}
