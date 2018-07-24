package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/go-redis/redis"
	"github.com/mengkzhaoyun/gostream/src/model"
)

// Receiver ,
type Receiver func(model.EventMessage)

type subscriber struct {
	receiver Receiver
}

// Publisher ,
type Publisher struct {
	sync.Mutex
	redis  *redis.Client
	topics map[string]*topic
}

// NewPubsub creates an in-memory publisher.
func NewPubsub(url string) Publisher {
	p := &Publisher{
		topics: make(map[string]*topic),
	}

	addr := url
	if strings.HasPrefix(url, "redis://") {
		addr = strings.TrimPrefix(url, "redis://")
	}
	p.redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong, err := p.redis.Ping().Result()
	fmt.Println(pong, err)

	return *p
}

// Publish ,
func (p *Publisher) Publish(c context.Context, dest string, message model.EventMessage) error {
	msgByte, err := json.Marshal(message)
	err = p.redis.Publish(dest, string(msgByte)).Err()
	return err
}

// Subscribe ,
func (p *Publisher) Subscribe(c context.Context, dest string, receiver Receiver) error {
	p.Lock()
	t, ok := p.topics[dest]
	if !ok {
		pubsub := p.redis.Subscribe(dest)
		msgC := pubsub.Channel()
		t = newTopic(dest, msgC)
		p.topics[dest] = t
	}
	p.Unlock()

	s := &subscriber{
		receiver: receiver,
	}
	t.subscribe(s)
	select {
	case <-c.Done():
	case <-t.done:
		t.unsubscribe(s)
	}
	return nil
}
