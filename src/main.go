package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mengkzhaoyun/gostream/src/conf"
	"github.com/mengkzhaoyun/gostream/src/router"
	"github.com/mengkzhaoyun/gostream/src/router/middleware"
	"github.com/mengkzhaoyun/gostream/src/sse"
	"github.com/mengkzhaoyun/gostream/src/version"

	"github.com/gin-gonic/contrib/ginrus"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var flags = []cli.Flag{
	cli.BoolFlag{
		EnvVar: "GOSTREAM_DEBUG",
		Name:   "debug",
		Usage:  "enable server debug mode",
	},
	cli.IntFlag{
		EnvVar: "GOSTREAM_SERVER_PORT",
		Name:   "server-port",
		Usage:  "http server port",
		Value:  80,
	},
	cli.StringFlag{
		EnvVar: "GOSTREAM_SERVER_PREFIX",
		Name:   "server-prefix",
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
	conf.Services.Pubsub = sse.NewPubsub()
	conf.Services.Prefix = c.String("server-prefix")

	handler := router.Load(
		ginrus.Ginrus(logrus.StandardLogger(), time.RFC3339, true),
		middleware.Version,
	)

	return http.ListenAndServe(
		fmt.Sprintf(":%d", c.Int("server-port")),
		handler,
	)
}

func before(c *cli.Context) error { return nil }
