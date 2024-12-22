package main

import (
	"log"
	"my-message-queue/queue"
)

func main() {
	conn := queue.ServerConnection{
		URL: "http://localhost:6060",
	}

	ok := conn.CreateQueue("test1")
	if !ok {
		log.Printf("failed creating queue")
		return
	}
	ok = conn.CreateQueue("test2")
	if !ok {
		log.Printf("failed creating queue")
		return
	}

	ok = conn.Publish("test1", "test msg")
	if !ok {
		log.Printf("failed publishing to queue")
		return
	}

	consumer := conn.Consume("test1")
	go func() {
		for message := range consumer {
			log.Print(message)
		}
	}()

	consumer2 := conn.Consume("test2")
	for message := range consumer2 {
		log.Print(message)
	}

}
