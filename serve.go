package bigbro

import (
	"github.com/gin-gonic/gin"
	"log"
)

func (l Logger) GinEndpoint(c *gin.Context) {
	ws, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go WsEvent(ws, l)
}
