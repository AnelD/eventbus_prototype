package bus

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type EventBus struct {
	// mutex to get locks for concurrency
	mutex sync.RWMutex
	// go doesn't have sets so map[T]bool is used
	subscribers map[string]map[*websocket.Conn]bool
}

func New() *EventBus {
	// constructor
	return &EventBus{
		subscribers: make(map[string]map[*websocket.Conn]bool),
	}
}

func (eb *EventBus) Subscribe(topic string, conn *websocket.Conn) {
	// get lock on eventbus to safely modify subscribers
	eb.mutex.Lock()
	// release lock when function is complete
	defer eb.mutex.Unlock()

	// check if topic already exists
	if eb.subscribers[topic] == nil {
		// if it doesn't create a map(set) of connections for it
		eb.subscribers[topic] = make(map[*websocket.Conn]bool)
	}
	// if it exists just add the connection
	eb.subscribers[topic][conn] = true
}

func (eb *EventBus) Publish(topic string, message Message, conn *websocket.Conn) {
	// get Read lock for subscribers
	eb.mutex.RLock()
	// release lock when complete
	defer eb.mutex.RUnlock()

	// send message to all subscribers
	for conn := range eb.subscribers[topic] {
		conn.WriteJSON(message)
	}
}

func (eb *EventBus) Unsubscribe(conn *websocket.Conn) {
	// get and release lock
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	// go through all topics
	for topic, conns := range eb.subscribers {
		// if conn is subscribed delete it
		if conns[conn] {
			delete(conns, conn)
			// if no connections are left for a topic delete it
			if len(conns) == 0 {
				delete(eb.subscribers, topic)
			}
		}
	}

}

func (eb *EventBus) LogAll() {
	for {
		time.Sleep(30 * time.Second)
		eb.mutex.RLock()
		log.Println("Current subscribers:", eb.subscribers)
		eb.mutex.RUnlock()
	}
}
