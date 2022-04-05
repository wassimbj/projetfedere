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

var clients = make(map[int]*websocket.Conn) // connected clients

type Message struct {
	Body   string `json:"msg"`
	SentTo int    `json:"sentTo"`
}

func (C) Start(res http.ResponseWriter, req *http.Request) {
	ws, err := upgrader.Upgrade(res, req, nil)

	session, _ := config.NewSession(req, res).Get("user")
	loggedInUser := session.Values
	myId := int(loggedInUser["id"].(int64))
	if err == nil {
		clients[int(loggedInUser["id"].(int64))] = ws
		fmt.Println("Connected: ", loggedInUser["id"])
	}
	// go (func() {
	for {
		defer ws.Close()
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("READ ERROR: ", err)
			delete(clients, myId)
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

		// send back the message to me
		if clients[myId] != nil {
			err1 := clients[myId].WriteJSON(msg)
			if err1 != nil {
				log.Printf("error: %v", err1)
				clients[myId].Close()
				delete(clients, myId)
			}
		} else {
			fmt.Println(myId, "Not Found !")
		}

		// send the message to the recipient
		if clients[msg.SentTo] != nil {
			err2 := clients[msg.SentTo].WriteJSON(msg)
			if err2 != nil {
				log.Printf("error: %v", err2)
				clients[msg.SentTo].Close()
				delete(clients, msg.SentTo)
			}
		} else {
			fmt.Println(msg.SentTo, "Not Found !")
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
