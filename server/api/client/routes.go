package client_routes

import (
	"net/http"
	"time"

	"pfserver/api/client/handlers"
	"pfserver/core"
	"pfserver/middlewares"

	"github.com/gorilla/mux"
	"github.com/wassimbj/gorl"
)

func ClientApiRoutes(router *mux.Router) {

	router.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		core.Respond(res, core.ResOpts{
			Status: 200,
			Msg:    "Hello 👋",
		})
	}).Methods("GET")

	//!################# auth endpoints #################
	// signup
	router.HandleFunc(
		"/signup",
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
		"/login",
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
		"/logout",
		middlewares.CheckAuth(handlers.Auth().Logout, true),
	).Methods("GET")

	// get logged-in user details
	router.HandleFunc(
		"/me",
		middlewares.CheckAuth(handlers.Auth().Me, true),
	).Methods("GET")

}
