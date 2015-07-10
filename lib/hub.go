package lib

import (
	// "log"
	"time"
)

type Hub struct {
	Connections map[*Connection]bool
	Broadcast   chan []byte
	Register    chan *Connection
	Unregister  chan *Connection
}

func NewHub() *Hub {

	return &Hub{
		Broadcast:   make(chan []byte),
		Register:    make(chan *Connection),
		Unregister:  make(chan *Connection),
		Connections: make(map[*Connection]bool),
	}

}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.Register:
			h.Connections[c] = true
		case c := <-h.Unregister:
			if _, ok := h.Connections[c]; ok {
				delete(h.Connections, c)
				close(c.send)
			}
		case m := <-h.Broadcast:
			// log.Println("m:", h.Connections)
			for c := range h.Connections {
				// log.Printf("%#v", c)
				select {
				case c.send <- m:
					// log.Println(m)
				default:
					delete(h.Connections, c)
					close(c.send)
				}
			}
		}
	}
}

func (h *Hub) Productmessage(w *Tips) {
	timer := time.NewTicker(5 * time.Second)
	for {
		<-timer.C
		m := w.GetMessage()
		h.Broadcast <- []byte(m)

	}

}
