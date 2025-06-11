package filewatcher

import (
	"encoding/json"
	"log"
	"time"

	"github.com/AnelD/eventbus/bus"
	"github.com/gorilla/websocket"
)

func PublishFileEventsWS(path string, shutdown chan struct{}) {
	var conn *websocket.Conn
	var err error
	for i := 0; i < 5; i++ {
		conn, _, err = websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
		if err == nil {
			break
		}
		log.Println("Dial failed, retrying in 500ms:", err)
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		log.Fatal("dial failed after retries:", err)
	}
	defer conn.Close()

	eventChan := make(chan string)

	go Watch(path, eventChan, shutdown)

	go publishEvents(conn, eventChan, shutdown)

	<-shutdown

}

func publishEvents(conn *websocket.Conn, eventChan chan string, shutdown chan struct{}) {
	for {
		select {
		case <-shutdown:
			return
		default:
			msg := bus.Message{Type: "publish", Topic: "file.upload", Data: <-eventChan}
			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				log.Println("Error marshaling message for log:", err)
			} else {
				log.Println("sending message: ", string(jsonMsg))
			}
			err = conn.WriteJSON(msg)
			if err != nil {
				log.Fatal("Error publishing file upload event:", err)
			}
		}
	}
}
