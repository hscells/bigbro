package main

import (
	"fmt"
	"github.com/buger/goterm"
	"github.com/gin-gonic/gin"
	"github.com/hscells/bigbro"
	"log"
	"time"
)

type server struct {
	l bigbro.Logger
}

func main() {
	t := time.Now().Format(time.Stamp)
	logger, err := bigbro.NewLogger(fmt.Sprintf("bigbrother_%s.log", t), bigbro.CSVFormatter{})
	if err != nil {
		log.Fatalln(err)
	}

	s := server{
		l: logger,
	}

	g := gin.Default()

	g.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	g.GET("/event", s.handleEvent)
	if goterm.Width() > 91 {
		fmt.Print(`

 @@@@@@@  @@@  @@@@@@@       @@@@@@@  @@@@@@@   @@@@@@  @@@@@@@ @@@  @@@ @@@@@@@@ @@@@@@@ 
 @@!  @@@ @@! !@@            @@!  @@@ @@!  @@@ @@!  @@@   @@!   @@!  @@@ @@!      @@!  @@@
 @!@!@!@  !!@ !@! @!@!@      @!@!@!@  @!@!!@!  @!@  !@!   @!!   @!@!@!@! @!!!:!   @!@!!@! 
 !!:  !!! !!: :!!   !!:      !!:  !!! !!: :!!  !!:  !!!   !!:   !!:  !!! !!:      !!: :!! 
 :: : ::  :    :: :: :       :: : ::   :   : :  : :. :     :     :   : : : :: :::  :   : :
                                                                    ...is always watching

 Harry Scells 2018

`)
	} else {
		fmt.Print(`Big Brother
...is always watching

Harry Scells 2018
`)
	}
	g.Run("0.0.0.0:1984")
}
