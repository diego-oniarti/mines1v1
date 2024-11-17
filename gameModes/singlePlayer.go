package gamemodes

import (
	"bytes"
	"encoding/binary"
	"log"
	"net/http"
	"strconv"

	"github.com/diego-oniarti/mines1v1/shared"
	"github.com/gin-gonic/gin"
	_ "github.com/gorilla/websocket"
)

type GameParams struct {
    width  uint16;
    height uint16;
    bombs  uint16;
    tempo  uint16;
    timed  bool;
}

var games_params map[string]GameParams;
func init() {
    games_params = make(map[string]GameParams);
}

func SinglePlayerWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
        log.Println("Cannot create websocket");
		c.String(http.StatusInternalServerError, "Impossibile creare WebSocket")
		return
	}
	defer conn.Close()

	// for {
	// 	messageType, message, err := conn.ReadMessage()
	// 	if err != nil {
	// 		break
	// 	}
	// 	// Echo di ritorno al client
	// 	err = conn.WriteMessage(messageType, message)
	// 	if err != nil {
	// 		break
	// 	}
	// }

	messageType, game_id, err := conn.ReadMessage()
	if err != nil || messageType!=1 {
        return
	}
    game_id_str := string(game_id[:])
    game_params := games_params[game_id_str]

    vals := []uint16{
        game_params.width,
        game_params.height,
        game_params.bombs,
        game_params.tempo,
    };
    log.Println("vals: ", vals)

    buffer := new(bytes.Buffer);
    binary.Write(buffer, binary.BigEndian, vals);
    err = conn.WriteMessage(2, buffer.Bytes())
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

    if timed == "off" {
        tempo = 0;
    }

    game_id := shared.RandomString(6, "");
    for {
        if _, ok := games_params[game_id]; !ok { break }
    }

    games_params[game_id] = GameParams{
    	width:  uint16(width),
    	height: uint16(height),
    	bombs:  uint16(bombs),
    	timed:  timed!="off",
    	tempo:  uint16(tempo),
    }

    shared.Render(c, http.StatusOK, "singlePlayer.html", gin.H{"game_id": game_id});
}
