package gamemodes

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Configura l'origine se necessario per la sicurezza (qui lo lasciamo aperto per testing)
		return true
	},
}
