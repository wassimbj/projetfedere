package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"pfserver/config"
	"pfserver/services"
	"pfserver/utils"
	"strconv"

	"github.com/gorilla/mux"
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
	Body   string `json:"msg"`
	SentTo int    `json:"sentTo"`
}

func (C) Start(res http.ResponseWriter, req *http.Request) {
	ws, err := upgrader.Upgrade(res, req, nil)
	if err == nil {
		clients[ws] = true
	}

	session, _ := config.NewSession(req, res).Get("user")
	loggedInUser := session.Values

	// go (func() {
	for {
		defer ws.Close()
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("READ ERROR: ", err)
			delete(clients, ws)
			return
		}

		createErr := services.Chat().Create(req.Context(), services.CreateData{
			SentFrom: int(loggedInUser["id"].(int64)),
			SentTo:   msg.SentTo,
			Msg:      msg.Body,
		})

		if createErr != nil {
			log.Printf("create msg error: %v", createErr)
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
	// })()
}

func (C) Messages(res http.ResponseWriter, req *http.Request) {
	chatTmpl, _ := template.ParseFiles(utils.TemplatePath("chat.html"))
	otherPeerId, _ := strconv.Atoi(mux.Vars(req)["userId"])
	otherPeerDetails, _ := services.User().GetUserData(req.Context(), services.GetUserBy{
		Id: int64(otherPeerId),
	})

	loggedInUser, _ := config.NewSession(req, res).GetUser()

	messages, err := services.Chat().Get(
		req.Context(),
		otherPeerId,
		int(loggedInUser["id"].(int64)),
	)

	fmt.Println("messages: ", messages)
	if err != nil {
		chatTmpl.Execute(res, nil)
		return
	}

	type ChatData struct {
		UserId    int
		Messages  []*services.Message
		OtherPeer services.UserData
	}

	data := ChatData{
		UserId:    int(loggedInUser["id"].(int64)),
		Messages:  messages,
		OtherPeer: otherPeerDetails,
	}
	chatTmpl.Execute(res, data)
}
