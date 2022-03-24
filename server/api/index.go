package routes

import (
	client_routes "pfserver/api/client"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	client_routes.ClientApiRoutes(router)

	return router

}
