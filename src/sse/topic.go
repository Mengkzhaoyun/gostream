package sse

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-redis/redis"
	"github.com/mengkzhaoyun/gostream/src/model"
)

type topic struct {
	sync.Mutex

	name string
	done chan bool
	subs map[*subscriber]struct{}
	msgc <-chan *redis.Message
}

func newTopic(dest string, msgC <-chan *redis.Message) *topic {
	t := &topic{
		name: dest,
		done: make(chan bool),
		subs: make(map[*subscriber]struct{}),
		msgc: msgC,
	}

	go t.listen()
	return t
}

func (t *topic) subscribe(s *subscriber) {
	t.Lock()
	t.subs[s] = struct{}{}
	t.Unlock()
}

func (t *topic) unsubscribe(s *subscriber) {
	t.Lock()
	delete(t.subs, s)
	t.Unlock()
}

func (t *topic) listen() {
	for {
		select {
		case <-t.done:
			fmt.Println("case <-t.done")
			return
		case msg, ok := <-t.msgc:
			fmt.Println("case <-t.docase msg, ok := <-t.msgc:")
			if ok {
				msgObj := new(model.EventMessage)
				err := json.Unmarshal([]byte(msg.Payload), &msgObj)
				if err == nil {
					for sub := range t.subs {
						sub.receiver(*msgObj)
					}
				}
			}
		}
	}
}

func (t *topic) close() {
	t.Lock()
	close(t.done)
	t.Unlock()
}
