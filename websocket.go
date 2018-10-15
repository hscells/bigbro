package bigbro

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

// Upgrader upgrades a web socket.
var Upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return false
	},
}

// WsEvent handles read events from the web socket.
func WsEvent(ws *websocket.Conn, l Logger) {

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
			var event Event
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
