// Copyright 2018 Drone.IO Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package router

import (
	"fmt"
	"net/http"

	"github.com/mengkzhaoyun/gostream/src/conf"
	"github.com/mengkzhaoyun/gostream/src/router/middleware/header"
	"github.com/mengkzhaoyun/gostream/src/server"

	"github.com/gin-gonic/gin"
)

// Load loads the router
func Load(middleware ...gin.HandlerFunc) http.Handler {

	e := gin.New()
	e.Use(gin.Recovery())

	e.Use(header.NoCache)
	e.Use(header.Options)
	e.Use(header.Secure)
	e.Use(middleware...)

	sse := e.Group(fmt.Sprintf("%s%s", conf.Services.Prefix, "/stream"))
	{
		sse.GET("/events", server.EventStreamSSE)
	}

	msg := e.Group(fmt.Sprintf("%s%s", conf.Services.Prefix, "/message"))
	{
		msg.POST("/events", server.EventStreamMSG)
	}

	e.GET(fmt.Sprintf("%s%s", conf.Services.Prefix, "/version"), server.Version)
	e.GET(fmt.Sprintf("%s%s", conf.Services.Prefix, "/healthz"), server.Health)

	return e
}
