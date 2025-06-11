package ws

import (
	"net/http"

	"github.com/AnelD/eventbus/bus"
)

func NewRouter(eb *bus.EventBus) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", HandleWS(eb))
	return mux
}
