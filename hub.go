package main

type Speakers struct {
	Present []string `json:"present"`
	Absent  []string `json:"absent"`
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan Speakers
	register   chan *Client
	unregister chan *Client
	speakers   Speakers
}

func newHub(speakers Speakers) *Hub {
	return &Hub{
		broadcast:  make(chan Speakers),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		speakers:   speakers,
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			client.send <- h.speakers
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.send)
				delete(h.clients, client)
			}
		case msg := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
