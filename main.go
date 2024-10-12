package main

import (
	"fmt"
	"sync"
	"time"
)

type Message struct {
	ID   int
	Body string
}

type Queue struct {
	messages []Message
	mutex    sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		messages: make([]Message, 0),
	}
}

func (q *Queue) Enqueue(msg Message) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.messages = append(q.messages, msg)
	fmt.Printf("Enqueued message: %d\n", msg.ID)
}

func (q *Queue) Dequeue() (Message, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	if len(q.messages) == 0 {
		return Message{}, false
	}
	msg := q.messages[0]
	q.messages = q.messages[1:]
	fmt.Printf("Dequeued message: %d\n", msg.ID)
	return msg, true
}

func producer(q *Queue, start, end int) {
	for i := start; i <= end; i++ {
		msg := Message{ID: i, Body: fmt.Sprintf("Message body %d", i)}
		q.Enqueue(msg)
		time.Sleep(time.Millisecond * 100)
	}
}

func consumer(q *Queue, id int) {
	for {
		if msg, ok := q.Dequeue(); ok {
			fmt.Printf("Consumer %d processed message: %d\n", id, msg.ID)
			time.Sleep(time.Millisecond * 300)
		} else {
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func main() {
	queue := NewQueue()

	// Start producers
	go producer(queue, 1, 5)
	go producer(queue, 6, 10)

	// Start consumers
	go consumer(queue, 1)
	go consumer(queue, 2)

	// Let the simulation run for a while
	time.Sleep(time.Second * 5)
}
