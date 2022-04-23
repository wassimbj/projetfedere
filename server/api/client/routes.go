package client_routes

import (
	"net/http"
	"path"
	"text/template"

	"pfserver/api/client/handlers"
	"pfserver/middlewares"
	"pfserver/utils"

	"github.com/gorilla/mux"
)

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }
// var clients = make(map[*websocket.Conn]bool) // connected clients
// var broadcast = make(chan Message)           // broadcast channel

func ClientApiRoutes(router *mux.Router) {

	// parse static files, like css, images...
	stylesPath := path.Join(utils.RootPath(), "ui/")
	fs := http.FileServer(http.Dir(stylesPath))
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", fs))

	loginTmpl, _ := template.ParseFiles(utils.TemplatePath("login.html"))
	signupTmpl, _ := template.ParseFiles(utils.TemplatePath("signup.html"))
	settingsTmpl, _ := template.ParseFiles(utils.TemplatePath("settings.html"))

	router.HandleFunc("/login", middlewares.CheckAuth(func(res http.ResponseWriter, req *http.Request) {
		loginTmpl.Execute(res, nil)
	}, false))

	router.HandleFunc("/signup", middlewares.CheckAuth(func(res http.ResponseWriter, req *http.Request) {
		signupTmpl.Execute(res, nil)
	}, false))

	router.HandleFunc("/settings", middlewares.CheckAuth(func(res http.ResponseWriter, req *http.Request) {
		settingsTmpl.Execute(res, nil)
	}, true))

	router.HandleFunc("/", middlewares.CheckAuth(handlers.Users().MembersList, true))

	router.HandleFunc("/chat/{userId}", middlewares.CheckAuth(handlers.Chat().Messages, true))

	//!################# chat endpoints #################
	router.HandleFunc("/ws", middlewares.CheckAuth(handlers.Chat().Start, true))

	//!################# auth endpoints #################
	// signup
	router.HandleFunc(
		"/api/signup",
		middlewares.CheckAuth(handlers.Auth().Signup, false),
	).Methods("POST")

	// login
	router.HandleFunc(
		"/api/login",
		middlewares.CheckAuth(handlers.Auth().Login, false),
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

	// get logged-in user details
	router.HandleFunc(
		"/api/settings",
		middlewares.CheckAuth(handlers.Users().Settings, true),
	).Methods("POST")

	router.HandleFunc(
		"/api/block",
		middlewares.CheckAuth(handlers.Users().BlockUnBlock, true),
	).Methods("POST")
}
