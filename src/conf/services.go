package conf

import "github.com/mengkzhaoyun/gostream/src/sse"

// Services is an evil global configuration that will be used as we transition /
// refactor the codebase to move away from storing these values in the Context.
var Services = struct {
	Pubsub sse.Publisher
	Prefix string
}{}
