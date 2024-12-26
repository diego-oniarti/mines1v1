package gamemodes

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/diego-oniarti/mines1v1/shared"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/gorilla/websocket"
)

func M1v1Ws(c *gin.Context) {
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

    var game *Game

    for {
        isFirstMove := true;

        remaining_time := time.Duration(game_params.tempo)*time.Millisecond
        var last_move time.Time

        for game==nil || game.state==Running {
            if !isFirstMove && game_params.timed {
                conn.SetReadDeadline(time.Now().Add(remaining_time))
            }
            messageType, move, err := conn.ReadMessage()
            if err != nil {
                if websocket.IsCloseError(err, 1001) {
                    return
                }
                changes := game.get_loosing_message()
                game.state=Lost
                send_changes(&changes, conn, game.state)
                return
            }
            if messageType!=2 { return } // messageType: 1=text; 2=binary
            x,y,flag := bytesToMove(move)

            if isFirstMove {
                if flag { continue }
                game = NewGame(game_params.width, game_params.height,
                game_params.bombs, game_params.tempo,
                x,y)
                isFirstMove=false
            }

            if flag {
                flagged, err := game.flag(x, y)
                if err==nil {
                    remaining_time = remaining_time - time.Now().Sub(last_move)
                    send_flagged(flagged, x, y, conn)
                }
            }else{
                changes, err := game.click(x, y)
                if err==nil {
                    send_changes(&changes, conn, game.state)
                    remaining_time = time.Duration(game_params.tempo)*time.Millisecond
                }
            }
            last_move = time.Now()
        }

        game=nil
        messageType, _, err := conn.ReadMessage()
        if err!=nil || messageType!=2 {
            return
        }
    }
}

func M1v1Page(c *gin.Context) {
    var width, height, bombs, tempo int;
    valid := true;
    tmp := func(name, default_v string) (int) {
        v,e := strconv.Atoi(c.DefaultQuery(name, default_v))
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

    // Generate game ids until you get a free one
    var game_id string
    for {
        game_id = shared.RandomString(6, "");
        if _, present := games_params[game_id]; !present { break }
    }

    games_params[game_id] = GameParams{
        width:  uint16(width),
        height: uint16(height),
        bombs:  uint16(bombs),
        timed:  timed!="off",
        tempo:  uint16(tempo),
    }

    shared.Render(c, http.StatusOK, "1v1.html", gin.H{"game_id": game_id});
}
