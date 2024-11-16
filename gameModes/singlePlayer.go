package gamemodes

import (
	"net/http"
	"strconv"

	"github.com/diego-oniarti/mines1v1/shared"
	"github.com/gin-gonic/gin"
	_ "github.com/gorilla/websocket"
)

func SinglePlayerWs(c *gin.Context) {
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

func SinglePlayerPage(c *gin.Context) {
    var width, height, bombs, tempo int;
    valid := true;
    tmp := func(name, def string) (int) {
        v,e := strconv.Atoi(c.DefaultQuery(name, def))
        if e!=nil || v<=0 { valid=false; return -1; }
        return v;
    }
    width = tmp("width" , "18")
    height = tmp("height", "14")
    bombs = tmp("bombs" , "40")
    timed := c.DefaultQuery("timed" , "off")
    tempo = tmp("tempo" , "3000")
    if !valid {
        c.Status(400);
        return;
    }

    _,_,_,_,_ = width,height,bombs,tempo,timed;
    
    shared.Render(c, http.StatusOK, "singlePlayer.html", nil);
}
