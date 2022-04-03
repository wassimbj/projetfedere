package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type C struct{}

func Chat() C {
	return C{}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var clients = make(map[*websocket.Conn]bool) // connected clients

type Message struct {
	Body string `json:"msg"`
}

func (C) Start(res http.ResponseWriter, req *http.Request) {
	ws, err := upgrader.Upgrade(res, req, nil)
	if err == nil {
		clients[ws] = true
	}

	go (func() {
		for {
			defer ws.Close()
			var msg Message
			err := ws.ReadJSON(&msg)
			if err != nil {
				log.Println("READ ERROR: ", err)
				delete(clients, ws)
				return
			}

			for client := range clients {
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	})()
}
