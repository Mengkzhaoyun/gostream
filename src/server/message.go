package server

import (
	"net/http"

	"github.com/Mengkzhaoyun/gostream/src/conf"
	"github.com/cncd/pubsub"
	"github.com/gin-gonic/gin"
)

// EventStreamMSG , Event
func EventStreamMSG(c *gin.Context) {
	in := new(pubsub.Message)
	if err := c.Bind(in); err != nil {
		c.String(http.StatusBadRequest, "Error parsing secret. %s", err)
		return
	}

	msg := &pubsub.Message{
		ID:     "",
		Data:   in.Data,
		Labels: in.Labels,
	}

	conf.Services.Pubsub.Publish(c, "topic/events", *msg)

	c.String(200, "Message Publish to topic/evnets : %s", msg.ID)
}
