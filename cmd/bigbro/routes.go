package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hscells/bigbro"
	"log"
	"net/http"
	"time"
)

// upgrader upgrades a web socket.
var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// wsEvent handles read events from the web socket.
func wsEvent(ws *websocket.Conn, l bigbro.Logger) {

	readWait := 1 * time.Millisecond
	readTicker := time.NewTicker(readWait)

	// defer closing of web socket
	defer func() {
		readTicker.Stop()
		ws.Close()
	}()

	for {
		select {
		case <-readTicker.C:
			var event bigbro.Event
			err := ws.ReadJSON(&event)
			if err != nil {
				log.Println(err)
				return
			}
			err = l.Log(event)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

// handleEvent handles an incoming request and attempts to upgrade it to a websocket.
func (s server) handleEvent(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go wsEvent(ws, s.l)
}
