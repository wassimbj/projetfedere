package main

import (
	"fmt"
	"log"

	routes "pfserver/api"

	"pfserver/config"
	"pfserver/core"
)

func main() {

	router := routes.Router()

	PORT := config.GetEnv("PORT")

	fmt.Printf("Server is running, http://localhost:%s", PORT)
	log.Fatal(core.StartServer(router, PORT))
}
