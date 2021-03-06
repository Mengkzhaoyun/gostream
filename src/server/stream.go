package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mengkzhaoyun/gostream/src/conf"
	"github.com/mengkzhaoyun/gostream/src/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// EventStreamSSE , Event
func EventStreamSSE(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	rw := c.Writer

	flusher, ok := rw.(http.Flusher)
	if !ok {
		c.String(500, "Streaming not supported")
		return
	}

	// ping the client
	io.WriteString(rw, ": ping\n\n")
	flusher.Flush()

	logrus.Debugf("user feed: connection opened")

	eventc := make(chan string)
	ctx, cancel := context.WithCancel(
		context.Background(),
	)

	defer func() {
		cancel()
		close(eventc)
		logrus.Debugf("user feed: connection closed")
	}()

	go func() {
		// TODO remove this from global config
		conf.Services.Pubsub.Subscribe(c, "topic/events", func(m model.EventMessage) {
			select {
			case <-ctx.Done():
				return
			default:
				eventc <- m.Data
			}
		})
		cancel()
	}()

	for {
		select {
		case <-rw.CloseNotify():
			return
		case <-ctx.Done():
			return
		case <-time.After(time.Second * 30):
			io.WriteString(rw, ": ping\n\n")
			flusher.Flush()
		case buf, ok := <-eventc:
			fmt.Println("buf, ok := <-eventc:")
			if ok {
				io.WriteString(rw, "data: ")
				rw.WriteString(buf)
				io.WriteString(rw, "\n\n")
				flusher.Flush()
			}
		}
	}
}
