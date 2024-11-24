package gamemodes

import (
	"bytes"
	"encoding/binary"
	"log"
	"net/http"
	"strconv"

	"github.com/diego-oniarti/mines1v1/shared"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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

func bytesToMove(b []byte) (uint16, uint16, bool) {
    x := binary.BigEndian.Uint16(b[0:2])
    y := binary.BigEndian.Uint16(b[2:4])
    flag := b[4]>0
    return x,y,flag
}
func SinglePlayerWs(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Println("Cannot create websocket");
        c.String(http.StatusInternalServerError, "Impossibile creare WebSocket")
        return
    }
    defer conn.Close()

    messageType, game_id, err := conn.ReadMessage()
    if err != nil || messageType!=1 { return }

    game_id_str := string(game_id[:])
    game_params, ok := games_params[game_id_str]
    if !ok { return }
    delete(games_params, game_id_str)


    err = conn.WriteMessage(2, arrToBuff([]uint16{
        game_params.width,
        game_params.height,
        game_params.bombs,
        game_params.tempo,
    }))

    // game := NewGame(game_params.width, game_params.height, game_params.bombs, game_params.tempo)
    var game *Game

    isFirstMove := true;
    for game==nil || game.state==Running {
        messageType, move, err := conn.ReadMessage()
        if err != nil || messageType!=2 { return }
        x,y,flag := bytesToMove(move)

        if isFirstMove {
            game = NewGame(game_params.width, game_params.height, game_params.bombs, game_params.tempo, x,y)
            isFirstMove=false
        }

        if flag {
            flagged, err := game.flag(x, y)
            if err==nil {
                send_flagged(flagged, x, y, conn)
            }
        }else{
            changes, err := game.click(x, y)
            if err==nil {
                send_cnahges(&changes, conn, game.state)
            }
        }
    }
}

func send_flagged(flagged bool, x uint16, y uint16, conn *websocket.Conn) {
    buffer := new(bytes.Buffer)
    var bits byte = 64
    if flagged {bits+=32}
    binary.Write(buffer, binary.BigEndian, bits)
    binary.Write(buffer, binary.BigEndian, x)
    binary.Write(buffer, binary.BigEndian, y)
    conn.WriteMessage(2, buffer.Bytes())
}

func send_cnahges(changes *[]CellaCoords, conn *websocket.Conn, state GameState) {
    buffer := new(bytes.Buffer)
    var state_bits byte = 0
    if state != Running {
        state_bits+=32
        if state==Won {
            state_bits+=16
        }
    }
    binary.Write(buffer, binary.BigEndian, state_bits)
    for i, change := range(*changes) {
        binary.Write(buffer, binary.BigEndian, change.x)
        binary.Write(buffer, binary.BigEndian, change.y)
        var hasNext uint8
        if i==len(*changes)-1 {hasNext=0} else {hasNext=8}
        binary.Write(buffer, binary.BigEndian, (change.cella.label<<4) + hasNext)
    }
    conn.WriteMessage(2, buffer.Bytes())
}

func arrToBuff(arr []uint16) []byte {
    buffer := new(bytes.Buffer)
    binary.Write(buffer, binary.BigEndian, arr)
    return buffer.Bytes()
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
