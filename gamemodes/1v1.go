package gamemodes

import (
	"log"
	"net/http"
	"time"

	"github.com/diego-oniarti/mines1v1/shared"
	"github.com/gin-gonic/gin"
	_ "github.com/gorilla/websocket"
)

func M1v1Ws(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Cannot create websocket")
		c.String(http.StatusInternalServerError, "Impossibile creare WebSocket")
		return
	}
	defer conn.Close()

	// Il client manda un messaggio con il game_id (stringa)
	messageType, game_id, err := conn.ReadMessage()
	if err != nil || messageType != 1 {
		log.Println(err)
		return
	}

	game_id_str := string(game_id[:])
	game_instance, ok := games[game_id_str]
	if !ok {
		return
	} // se il game_id non esiste esci
	game_params := game_instance.params

	// Manda i parametri al client
	is_g1 := game_instance.g1 == nil

	var player_number uint16 = 0
	if !is_g1 {
		player_number = 1
	}
	err = conn.WriteMessage(2, arrToBuff([]uint16{
		game_params.width,
		game_params.height,
		game_params.bombs,
		game_params.tempo,
		player_number,
	}))

	to_other_chn := game_instance.a_to_b
	from_other_chn := game_instance.b_to_a
	other_conn := game_instance.g2 // Qui è nil
	if is_g1 {
		game_instance.g1 = conn
		defer delete(games, game_id_str) // Lascia chiudere solo il primo. Così se il secondo non arriva mai chiude comunque

		<-from_other_chn // Aspetta si sia connesso G2
		other_conn = game_instance.g2
	} else {
		game_instance.g2 = conn
		to_other_chn = game_instance.b_to_a
		from_other_chn = game_instance.a_to_b
		other_conn = game_instance.g1

		to_other_chn <- 1 // avvisa che si è connesso G2
	}

	var game *Game
	var timer <-chan time.Time
	isFirstMove := is_g1   // g2 non ha mai la prima mossa
	waiting_other := false // G2 inizia aspettando l'altro

	if !is_g1 { // Aspetta la prima mossa dal g1
		<-from_other_chn
		game = game_instance.game // Settato da g1
	}

	move_chn := make(chan []byte, 1)
	message_type_chn := make(chan int, 1)
	error_chn := make(chan error, 1)

	go func() {
		for {
			messageType, move, err := conn.ReadMessage()
			if err != nil {
				error_chn <- err
				return
			}
			message_type_chn <- messageType
			move_chn <- move
		}
	}()

	for {
		var move []byte
		var messageType int
		var err error
		select {
		case move = <-move_chn:
			messageType = <-message_type_chn
			if waiting_other {
				continue
			} // Scarta i messaggi mandati mentre in attesa
			if messageType != 1 && messageType != 2 {
				return
			}
			x, y, flag := bytesToMove(move)

			if isFirstMove {
				if flag {
					continue
				} // Ignora la prima mossa se è una flag

				// G1 crea il game
				if is_g1 {
					game = NewGame(game_params.width, game_params.height,
						game_params.bombs, game_params.tempo,
						x, y)
					game_instance.game = game
					timer = get_timer(&game_params)
				}
				isFirstMove = false
			}

			if flag {
				flagged, err := game.flag(x, y)
				if err != nil {
					log.Println(err)
					continue
				} else {
					send_flagged(flagged, x, y, conn, false)
					send_flagged(flagged, x, y, other_conn, true)
				}
			} else {
				changes, err := game.click(x, y)
				if err != nil {
					log.Println(err)
					continue
				}
				send_changes(&changes, conn, game.state, false)
				send_changes(&changes, other_conn, game.state, true)
				to_other_chn <- 1
				waiting_other = true // Dopo aver fatto la mossa inizia ad aspettare l'altro
			}
		case err = <-error_chn:
			log.Println(err)
			return
		case <-timer:
			if waiting_other {
				continue
			} // Scarta i messaggi mandati mentre in attesa
			changes := game.get_loosing_message()
			game.state = Lost
			send_changes(&changes, conn, game.state, false)
			send_changes(&changes, other_conn, game.state, true)
		case <-from_other_chn:
			waiting_other = false
			timer = get_timer(&game_params)
		}
	}
}

func M1v1Page(c *gin.Context) {
	game_id := c.Query("game_id")
	if game_id == "" {
		c.Status(400)
		return
	}
	// controlla che il game_id sia valido. Previene reflect injection
	if _, present := games[game_id]; !present {
		c.Status(400)
		return
	}
	shared.Render(c, http.StatusOK, "1v1.html", gin.H{"game_id": game_id})
}
