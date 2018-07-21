package rest

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Mengkzhaoyun/gostream/src/model"
	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

// SSEHandler , Server Side Event
type SSEHandler struct {
	// Events are pushed to this channel by the main events-gathering routine
	Notifier chan model.EventMessage

	// New client connections
	newClients chan model.EventClient

	// Closed client connections
	closingClients chan model.EventClient

	// Client connections registry
	clients map[*model.EventClient]bool
}

// NewSSEHandler , Create Users Handler
func NewSSEHandler(prefix string) (http.Handler, error) {
	handler := &SSEHandler{
		Notifier:       make(chan model.EventMessage),
		newClients:     make(chan model.EventClient),
		closingClients: make(chan model.EventClient),
		clients:        make(map[*model.EventClient]bool),
	}

	wsContainer := restful.NewContainer()
	wsContainer.EnableContentEncoding(true)

	apiWs := new(restful.WebService)

	apiWs.Path(prefix).
		Consumes(restful.MIME_JSON, "text/event-stream").
		Produces(restful.MIME_JSON, "text/event-stream")
	wsContainer.Add(apiWs)

	apiTags := []string{strings.TrimSuffix(prefix, "/")}

	apiWs.Route(apiWs.GET("/").To(handler.test).
		// docs
		Doc("get a sample event message").
		Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Writes(model.EventMessage{}).
		Returns(200, "OK", model.EventMessage{}))

	apiWs.Route(apiWs.GET("/test").To(handler.test).
		// docs
		Doc("get a sample event message").
		Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Writes(model.EventMessage{}).
		Returns(200, "OK", model.EventMessage{}))

	apiWs.Route(apiWs.GET("/stream").To(handler.stream).
		// docs
		Doc("get a event stream").
		Metadata(restfulspec.KeyOpenAPITags, apiTags))

	apiWs.Route(apiWs.POST("/event").To(handler.event).
		// docs
		Doc("post event message to server").
		Metadata(restfulspec.KeyOpenAPITags, apiTags).
		Reads(model.EventMessage{})) // from the request

	return wsContainer, nil
}

// GET http://localhost/{prefix}/sse/stream
//
func (handler SSEHandler) stream(request *restful.Request, response *restful.Response) {
	rw := response.ResponseWriter
	flusher, ok := rw.(http.Flusher)

	if !ok {
		response.WriteError(http.StatusInternalServerError, fmt.Errorf("%s", "Streaming unsupported!"))
		return
	}

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	// Each connection registers its own message channel with the Broker's connections registry
	messageClient := model.EventClient{
		ID:      "mengkzhaoyun@gmail.com",
		Message: make(chan model.EventMessage),
	}

	// Signal the broker that we have a new connection
	handler.newClients <- messageClient

	// Remove this client from the map of connected clients
	// when this handler exits.
	defer func() {
		handler.closingClients <- messageClient
	}()

	// Listen to connection close and un-register messageChan
	notify := rw.(http.CloseNotifier).CloseNotify()

	go func() {
		<-notify
		handler.closingClients <- messageClient
	}()

	for {

		// Write to the ResponseWriter
		// Server Sent Events compatible
		fmt.Fprintf(rw, "data: %s\n\n", <-messageClient.Message)

		// Flush the data immediatly instead of buffering it for later.
		flusher.Flush()
	}
}

// POST http://localhost/{prefix}/sse/event
//
func (handler SSEHandler) event(request *restful.Request, response *restful.Response) {
	msg := new(model.EventMessage)
	err := request.ReadEntity(&msg)
	if err == nil {
		handler.Notifier <- *msg
		response.WriteEntity("success")
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}
}

func (handler *SSEHandler) listen() {
	for {
		select {
		case s := <-handler.newClients:

			// A new client has connected.
			// Register their message channel
			handler.clients[&s] = true
			log.Printf("Client added. %d registered clients", len(handler.clients))
		case s := <-handler.closingClients:

			// A client has dettached and we want to
			// stop sending them messages.
			delete(handler.clients, &s)
			log.Printf("Removed client. %d registered clients", len(handler.clients))
		case event := <-handler.Notifier:

			// We got a new event from the outside!
			// Send event to all connected clients
			for client := range handler.clients {
				client.Message <- event
			}
		}
	}

}

// POST http://localhost/{prefix}/sse/event
//
func (handler SSEHandler) test(request *restful.Request, response *restful.Response) {
	msg := &model.EventMessage{
		ID: "12po85ss",
		Labels: map[string]string{
			"user":    "shucheng@spacesystech.com",
			"channel": "drone",
		},
	}
	response.WriteEntity(msg)
}
