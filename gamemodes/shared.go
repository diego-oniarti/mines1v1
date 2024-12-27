package gamemodes

import (
    "bytes"
    "encoding/binary"
    "net/http"

    "github.com/gorilla/websocket"
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

func send_flagged(flagged bool, x uint16, y uint16, conn *websocket.Conn) {
    buffer := new(bytes.Buffer)
    var bits byte = 64
    if flagged {bits+=32}
    binary.Write(buffer, binary.BigEndian, bits)
    binary.Write(buffer, binary.BigEndian, x)
    binary.Write(buffer, binary.BigEndian, y)
    conn.WriteMessage(2, buffer.Bytes())
}

func arrToBuff(arr []uint16) []byte {
    buffer := new(bytes.Buffer)
    binary.Write(buffer, binary.BigEndian, arr)
    return buffer.Bytes()
}

func send_changes(changes *[]CellaCoords, conn *websocket.Conn, state GameState) {
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

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // Configura l'origine se necessario per la sicurezza (qui lo lasciamo aperto per testing)
        return true
    },
}
