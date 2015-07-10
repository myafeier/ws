package lib

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type Connection struct {
	ws   *websocket.Conn
	send chan []byte
	h    *Hub
}

func (c *Connection) writer() {
	defer func() {
		c.ws.Close()
		// log.Println("writer close")
	}()
	for {
		message := <-c.send
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err)
			break
		}
	}
	// for message := range c.send {
	// 	err := c.ws.WriteMessage(websocket.TextMessage, message)
	// 	if err != nil {
	// 		log.Println(err)
	// 		break
	// 	}
	// }
	// log.Println("writer is over")
}
func (c *Connection) reader() {
	defer func() {
		c.ws.Close()
		// log.Println("reader close")
	}()
	timer := time.NewTicker(10 * time.Second)
	for {
		<-timer.C
		_, _, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		//c.h.Broadcast <- message
	}
}

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

type WsHandler struct {
	H *Hub
}

func (wsh WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}
	c := &Connection{send: make(chan []byte, 0), ws: ws, h: wsh.H}
	c.h.Register <- c
	defer func() {
		c.h.Unregister <- c
	}()

	go c.writer()
	c.reader()
}
