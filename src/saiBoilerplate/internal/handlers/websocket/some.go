package websocket

import (
	"log"
	"net/http"
)

func (s Server) handle(w http.ResponseWriter, r *http.Request) {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
}
