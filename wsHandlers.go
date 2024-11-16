package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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

func wsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "Impossibile creare WebSocket")
		return
	}
	defer conn.Close()

	// Loop per gestire i messaggi WebSocket
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Echo di ritorno al client
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			break
		}
	}
}
