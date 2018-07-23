package sse

import (
	"context"
	"sync"

	"github.com/cncd/pubsub"
)

type subscriber struct {
	receiver pubsub.Receiver
}

type publisher struct {
	sync.Mutex

	topics map[string]*topic
}

// NewPubsub creates an in-memory publisher.
func NewPubsub() pubsub.Publisher {
	p := &publisher{
		topics: make(map[string]*topic),
	}

	return p
}

func (p *publisher) Create(c context.Context, dest string) error {
	p.Lock()
	t, ok := p.topics[dest]
	if !ok {
		t = newTopic(dest)
		p.topics[dest] = t
	}
	p.Unlock()
	return nil
}

func (p *publisher) Publish(c context.Context, dest string, message pubsub.Message) error {
	p.Lock()
	t, ok := p.topics[dest]
	p.Unlock()
	if !ok {
		t = newTopic(dest)
		p.topics[dest] = t
	}
	t.publish(message)
	return nil
}

func (p *publisher) Subscribe(c context.Context, dest string, receiver pubsub.Receiver) error {
	p.Lock()
	t, ok := p.topics[dest]
	p.Unlock()
	if !ok {
		t = newTopic(dest)
		p.topics[dest] = t
	}
	s := &subscriber{
		receiver: receiver,
	}
	t.subscribe(s)
	select {
	case <-c.Done():
	case <-t.done:
	}
	t.unsubscribe(s)
	return nil
}

func (p *publisher) Remove(c context.Context, dest string) error {
	p.Lock()
	t, ok := p.topics[dest]
	if ok {
		delete(p.topics, dest)
		t.close()
	}
	p.Unlock()
	return nil
}
