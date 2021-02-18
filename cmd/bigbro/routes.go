package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hscells/bigbro"
	"log"
)

// handleEvent handles an incoming request and attempts to upgrade it to a websocket.
func (s server) handleEvent(c *gin.Context) {
	ws, err := bigbro.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go bigbro.WsEvent(ws, s.l)
}

// handleEvent handles an incoming request and attempts to upgrade it to a websocket.
func (s server) handleCapture(c *gin.Context) {
	ws, err := bigbro.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	go bigbro.WsRecord(ws)
}
