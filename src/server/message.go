package server

import (
	"net/http"

	"github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"

	"github.com/mengkzhaoyun/gostream/src/conf"
	"github.com/mengkzhaoyun/gostream/src/model"
)

// EventStreamMSG , Event
func EventStreamMSG(c *gin.Context) {
	in := new(model.EventMessage)
	if err := c.Bind(in); err != nil {
		c.String(http.StatusBadRequest, "Error parsing message. %s", err)
		return
	}

	msg := &model.EventMessage{
		ID:     uuid.Must(uuid.NewV4()).String(),
		Data:   in.Data,
		Labels: in.Labels,
	}

	conf.Services.Pubsub.Publish(c, "topic/events", *msg)

	c.String(200, "Message Publish to topic/events : %s", msg.ID)
}
