package main

import (
	"github.com/amitde69/my-message-queue/queue"
)

func main() {
	q := queue.NewQueue("test")
	q.Publish("testmsg 1")
	q.Publish("testmsg 2")
	q.Publish("testmsg 3")
	q.ListMessages()
	q.Consume()
	q.ListMessages()
	q.Consume()
	q.ListMessages()
	q.Consume()
	q.ListMessages()

}
