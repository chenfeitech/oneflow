package server

import (
	"github.com/gorilla/websocket"
	"utils/helper"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// serverWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	// if r.Method != "GET" {
	// 	http.Error(w, "Method not allowed", 405)
	// 	return
	// }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := helper.NewConnection(ws)
	helper.WsHub.Register(c)
	go c.WritePump()
	c.ReadPump()
}
