package bigbro

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

// upgrader upgrades a web socket.
var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,

}

func (l Logger) GinEndpoint(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	go WsEvent(ws, l)
}
