package server

import (
	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {
	router := mux.NewRouter()
	Initmongo()
	router.HandleFunc("/api/userlist", GetAllUser).Methods("GET", "OPTIONS")
	return router

}