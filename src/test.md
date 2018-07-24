package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "k8s.spacecig.com:6379",
	})
	defer client.Close()

	// This part causes the program to hang at ReceiveMessage()
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	pubsub := client.Subscribe("mychannel1")
	defer pubsub.Close()

	go func() {

		<-time.After(time.Second * 3)

		if err := client.Publish("mychannel1", "hello1").Err(); err != nil {
			panic(err)
		}

		<-time.After(time.Second * 3)

		if err := client.Publish("mychannel1", "hello2").Err(); err != nil {
			panic(err)
		}
	}()

	msgC := pubsub.Channel()

	for {
		select {
		case <-time.After(time.Second * 10):
			return
		case msg, ok := <-msgC:
			if ok {
				fmt.Println(msg.Channel, msg.Payload)
			}
		}
	}
}
