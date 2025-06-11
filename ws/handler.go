package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AnelD/eventbus/bus"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true // Allow all connections
}}

func HandleWS(bus *bus.EventBus) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Println(writer, "WebSocket upgrade failed", http.StatusInternalServerError)
			return
		}

		go handleClient(conn, bus)
	}
}

func handleClient(conn *websocket.Conn, eb *bus.EventBus) {

	defer func() {
		eb.Unsubscribe(conn)
		conn.Close()
	}()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		if messageType != websocket.TextMessage {
			log.Println("Unsupported message type")
			continue
		}

		var parsed bus.Message
		if err := json.Unmarshal(msg, &parsed); err != nil {
			log.Println("invalid JSON:", err)
			continue
		}

		switch parsed.Type {
		case "subscribe":
			eb.Subscribe(parsed.Topic, conn)
		case "publish":
			eb.Publish(parsed.Topic, parsed, conn)
		default:
			log.Println("Received bad message type:", parsed.Type)
		}
	}
}
