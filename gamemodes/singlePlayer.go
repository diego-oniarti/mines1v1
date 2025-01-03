package gamemodes

import (
	"log"
	"net/http"
	"time"

	"github.com/diego-oniarti/mines1v1/shared"
	"github.com/gin-gonic/gin"
	_ "github.com/gorilla/websocket"
)

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
    game_instance, ok := games[game_id_str]
    if !ok { return }
    game_params := game_instance.params
    defer delete(games, game_id_str)

    err = conn.WriteMessage(2, arrToBuff([]uint16{
        game_params.width,
        game_params.height,
        game_params.bombs,
        game_params.tempo,
    }))

    var game *Game
    var timer <-chan time.Time

    // loop esterno. Ogni ciclo è una partita
    for {
        isFirstMove := true;

        // loop interno. Ogni ciclo è una mossa nella partita
        for game==nil || game.state==Running {
            move_chn := make(chan []byte, 1)
            message_type_chn := make(chan int, 1)
            error_chn := make(chan error, 1)

            go func() {
                messageType, move, err := conn.ReadMessage()
                if err!=nil {
                    error_chn <- err
                    return
                }
                message_type_chn <- messageType
                move_chn <- move
            }()

            var move []byte
            var messageType int
            var err error

            if isFirstMove || !game_params.timed {
                select {
                case move = <-move_chn:
                    messageType = <-message_type_chn
                case err = <-error_chn:
                }
            }else{
                select {
                case move = <-move_chn:
                    messageType = <-message_type_chn
                case err = <-error_chn:
                case <-timer:
                    changes := game.get_loosing_message()
                    game.state=Lost
                    send_changes(&changes, conn, game.state, false)
                    move = <-move_chn
                    messageType = <-message_type_chn
                }
            }

            if err != nil {
                log.Println(err)
                return
            }
            if messageType!=2 && messageType!=1 { 
                log.Println("BBB");
                log.Println(messageType)
                return
            } // messageType: 1=text; 2=binary
            if messageType==1 {
                s := string(move[:])
                log.Println(s)
                if (s=="replay") {
                    break
                }else{
                    return
                }
            }

            x,y,flag := bytesToMove(move)

            if isFirstMove {
                if flag { continue }
                game = NewGame(game_params.width, game_params.height,
                game_params.bombs, game_params.tempo,
                x,y)
                timer = time.After(time.Duration(game_params.tempo)*time.Millisecond)
                isFirstMove=false
            }

            if flag {
                flagged, err := game.flag(x, y)
                if err==nil {
                    send_flagged(flagged, x, y, conn,false)
                }
            }else{
                timer = time.After(time.Duration(game_params.tempo)*time.Millisecond)
                changes, err := game.click(x, y)
                if err==nil {
                    send_changes(&changes, conn, game.state, false)
                }
            }
        }
        game=nil
        timer=nil
    }
}

func SinglePlayerPage(c *gin.Context) {
    game_id := c.Query("game_id")
    if game_id=="" {
        c.Status(400)
        return
    }
    if _, present := games[game_id]; !present {
        c.Status(400)
        return
    }
    shared.Render(c, http.StatusOK, "singlePlayer.html", gin.H{"game_id": game_id});
}

func get_timer(params *GameParams) <- chan time.Time{
    if !params.timed {
        return nil
    }
    return time.After(time.Duration(params.tempo) * time.Millisecond)
}
