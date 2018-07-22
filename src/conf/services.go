package conf

import (
	"github.com/cncd/pubsub"
)

// Services is an evil global configuration that will be used as we transition /
// refactor the codebase to move away from storing these values in the Context.
var Services = struct {
	Pubsub pubsub.Publisher
	Prefix string
}{}
