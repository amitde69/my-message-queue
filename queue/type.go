package queue

import (
	"fmt"
	"sync"

	"golang.org/x/exp/slices"
)

type Queue struct {
	Name     string
	Messages []Message
	Mutex    sync.Mutex
}

type Message struct {
	Ack     bool `default:"false"`
	Payload string
}

func NewQueue(name string) *Queue {
	return &Queue{Name: name}
}
func (q *Queue) Publish(m string) bool {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	q.Messages = append(q.Messages, Message{Payload: m})
	return true
}

func (q *Queue) GetLen() int {
	count := 0
	for _, messsage := range q.Messages {
		if !messsage.Ack {
			count++
		}
	}
	return count
}

func (q *Queue) ListMessages() {
	if q.GetLen() == 0 {
		fmt.Println("The queue is empty.")
		return
	}

	fmt.Printf("Messages in queue '%s':\n", q.Name)
	reverseMessages := append([]Message(nil), q.Messages...)
	slices.Reverse(reverseMessages)
	for i, msg := range reverseMessages {
		if !msg.Ack {
			fmt.Printf("[%d] %s\n", i+1, msg.Payload)
		}
	}
}
func (q *Queue) Consume() Message {
	var firstUnacked Message
	var firstIndex int
	q.Mutex.Lock()
	defer q.Mutex.Unlock()
	if q.GetLen() < 1 {
		return Message{}
	}
	for index, messsage := range q.Messages {
		if !messsage.Ack {
			firstUnacked = messsage
			firstIndex = index
			break
		}
	}
	q.Messages[firstIndex].Ack = true

	return firstUnacked
}
