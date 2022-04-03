package client_routes

import (
	"fmt"
	"net/http"
	"path"
	"text/template"
	"time"

	"pfserver/api/client/handlers"
	"pfserver/middlewares"
	"pfserver/utils"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/wassimbj/gorl"
)

func ClientApiRoutes(router *mux.Router) {

	// parse static files, like css, images...
	stylesPath := path.Join(utils.RootPath(), "ui/")
	fs := http.FileServer(http.Dir(stylesPath))
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", fs))

	loginTmpl, _ := template.ParseFiles(utils.TemplatePath("login.html"))
	signupTmpl, _ := template.ParseFiles(utils.TemplatePath("signup.html"))
	homeTmpl, _ := template.ParseFiles(utils.TemplatePath("home.html"))
	chatTmpl, _ := template.ParseFiles(utils.TemplatePath("chat.html"))

	router.HandleFunc("/login", middlewares.CheckAuth(func(res http.ResponseWriter, req *http.Request) {
		loginTmpl.Execute(res, nil)
	}, false))

	router.HandleFunc("/signup", middlewares.CheckAuth(func(res http.ResponseWriter, req *http.Request) {
		signupTmpl.Execute(res, nil)
	}, false))

	router.HandleFunc("/", middlewares.CheckAuth(func(res http.ResponseWriter, req *http.Request) {
		homeTmpl.Execute(res, nil)
	}, true))

	type Msg struct {
		Msg       string
		SentFrom  int
		SentTo    int
		CreatedAt time.Time
	}
	type ChatData struct {
		UserId   int
		Messages []Msg
	}

	router.HandleFunc("/chat", func(res http.ResponseWriter, req *http.Request) {
		data := ChatData{
			UserId: 10,
			Messages: []Msg{
				{Msg: "Hello", SentFrom: 10, SentTo: 100, CreatedAt: time.Now()},
				{Msg: "Hello Back !!", SentFrom: 100, SentTo: 10, CreatedAt: time.Now()},
				{Msg: "How are you ?", SentFrom: 100, SentTo: 10, CreatedAt: time.Now()},
				{Msg: "Fine thanks :)", SentFrom: 10, SentTo: 100, CreatedAt: time.Now()},
			},
		}
		chatTmpl.Execute(res, data)
	})

	//!################# chat endpoints #################
	// Configure the upgrader
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	router.HandleFunc("/ws", func(res http.ResponseWriter, req *http.Request) {
		conn, err1 := upgrader.Upgrade(res, req, nil) // error ignored for sake of simplicity
		fmt.Println(err1)

		// for {
		// Read message from browser
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

		// Write message back to browser
		if err = conn.WriteMessage(msgType, msg); err != nil {
			return
		}
		// }
	})

	//!################# auth endpoints #################
	// signup
	router.HandleFunc(
		"/api/signup",
		middlewares.RateLimit(
			middlewares.CheckAuth(handlers.Auth().Signup, false),
			gorl.RLOpts{
				Attempts:      5,
				Prefix:        "signup",
				BlockDuration: time.Hour * 24, // block for one day
				Duration:      time.Hour * 24,
			},
		),
	).Methods("POST")

	// login
	router.HandleFunc(
		"/api/login",
		middlewares.RateLimit(
			middlewares.CheckAuth(handlers.Auth().Login, false),
			gorl.RLOpts{
				Attempts:      10,
				Prefix:        "login",
				BlockDuration: time.Hour,
				Duration:      time.Minute, // 10 attempts per hour
			},
		),
	).Methods("POST")

	// logout
	router.HandleFunc(
		"/api/logout",
		middlewares.CheckAuth(handlers.Auth().Logout, true),
	).Methods("GET")

	// get logged-in user details
	router.HandleFunc(
		"/api/me",
		middlewares.CheckAuth(handlers.Auth().Me, true),
	).Methods("GET")

}
