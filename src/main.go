package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Mengkzhaoyun/gostream/src/rest"
	"github.com/Mengkzhaoyun/gostream/src/version"

	"github.com/urfave/cli"
)

var flags = []cli.Flag{
	cli.IntFlag{
		Name:  "port",
		Usage: "http server port",
		Value: 80,
	},
	cli.StringFlag{
		EnvVar: "SSE_PREFIX",
		Name:   "prefix",
		Usage:  "http server prefix (/sse)",
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "gostream"
	app.Version = version.Version.String()
	app.Usage = "go stream events"
	app.Action = server
	app.Flags = flags
	app.Before = before

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func server(c *cli.Context) error {
	if len(c.String("prefix")) > 0 {
		log.Printf("Using HTTP prefix: %s", c.String("prefix"))
	}

	// /sse
	sseURL := fmt.Sprintf("%s%s", c.String("prefix"), "/sse/")
	sseHandler, _ := rest.NewAdminHandler(sseURL)
	http.Handle(sseURL, sseHandler)

	log.Printf("Using HTTP port: %d", c.Int("port"))
	http.ListenAndServe(fmt.Sprintf(":%d", c.Int("port")), nil)
	return nil
}

func before(c *cli.Context) error { return nil }
