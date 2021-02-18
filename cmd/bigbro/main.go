package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/buger/goterm"
	"github.com/gin-gonic/gin"
	"github.com/hscells/autotls"
	"github.com/hscells/bigbro"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net/http"
	"os"
)

type server struct {
	l bigbro.Logger
}

type args struct {
	Format       string   `help:"how events should be formatted and written" arg:"positional"`
	Filename     string   `help:"filename to output logs to"`
	Index        string   `help:"index for Elasticsearch to use"`
	V            string   `help:"version for Elasticsearch event type"`
	URL          string   `help:"URL for Elasticsearch"`
	CheckOrigin  bool     `help:"enable or disable same-origin requests"`
	TLS          bool     `help:"listen and server over secure channel"`
	TLSWhitelist []string `help:"whitelisted URLs for TLS"`
}

func (args) Version() string {
	return "20.Dec.2018"
}

func (args) Description() string {
	return "server for receiving and processing bigbro events"
}

func main() {
	var (
		args   args
		err    error
		logger bigbro.Logger
	)
	p := arg.MustParse(&args)
	if args.Format == "elasticsearch" && args.V == "" && args.Index == "" {
		p.Fail("you must provide one of --index and --v")
	}
	if args.Format == "csv" && args.Filename == "" {
		p.Fail("you must provide one of --filename")
	}

	switch args.Format {
	case "csv":
		logger, err = bigbro.NewCSVLogger(args.Filename)
		if err != nil {
			log.Fatalln(err)
		}
	case "elasticsearch":
		logger, err = bigbro.NewElasticsearchLogger(args.Index, args.V, args.URL)
		if err != nil {
			log.Fatalln(err)
		}
	default:
		log.Fatalf("%s is not a valid log format\n", args.Format)
	}

	fmt.Printf("checking origin? %v\n", !args.CheckOrigin)
	bigbro.Upgrader.CheckOrigin = func(r *http.Request) bool {
		return !args.CheckOrigin
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

	port := os.Getenv("BIGBRO_PORT")
	if len(port) == 0 {
		port = "1984"
	}

	g.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain", []byte(fmt.Sprintf("bigbro version %s", args.Version())))
	})
	g.GET("/event", s.handleEvent)
	g.GET("/capture", s.handleCapture)
	if goterm.Width() > 91 {
		fmt.Printf(`

 @@@@@@@  @@@  @@@@@@@       @@@@@@@  @@@@@@@   @@@@@@  @@@@@@@ @@@  @@@ @@@@@@@@ @@@@@@@ 
 @@!  @@@ @@! !@@            @@!  @@@ @@!  @@@ @@!  @@@   @@!   @@!  @@@ @@!      @@!  @@@
 @!@!@!@  !!@ !@! @!@!@      @!@!@!@  @!@!!@!  @!@  !@!   @!!   @!@!@!@! @!!!:!   @!@!!@! 
 !!:  !!! !!: :!!   !!:      !!:  !!! !!: :!!  !!:  !!!   !!:   !!:  !!! !!:      !!: :!! 
 :: : ::  :    :: :: :       :: : ::   :   : :  : :. :     :     :   : : : :: :::  :   : :
                                                                    ...is always watching

 Harry Scells 2018
 version %s

`, args.Version())
	} else {
		fmt.Printf(`Big Brother
...is always watching

Harry Scells 2018
version %s
`, args.Version())
	}

	if args.TLS {
		m := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(args.TLSWhitelist...),
			Cache:      autocert.DirCache(".tls_cache"),
		}
		log.Fatal(autotls.RunWithManager(g, &m))
	} else {
		fmt.Printf("port: %s\n", port)
		g.Run(fmt.Sprintf("0.0.0.0:%s", port))
	}
}
