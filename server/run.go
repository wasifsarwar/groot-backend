package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

//Run the app
func Run() {
	r := Router()
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	fmt.Println()
	fmt.Println("Starting server on the port:", port, " ...")
	log.Fatal(http.ListenAndServe(":"+port, r))
	fmt.Println("this is heroku's port:", port)

}
