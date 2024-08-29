package router

import (
	"employeecurd/handlers"

	"github.com/gorilla/mux"
)

// NewRouter initializes a new router with the provided UserHandler interface
func NewRouter(userHandler handlers.UserHandler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	
	// Use the methods from the userHandler instance
	router.HandleFunc("/assignment/user", userHandler.GetUser).Methods("GET")
	router.HandleFunc("/assignment/user", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/assignment/user", userHandler.UpdateUser).Methods("PATCH")

	return router
}

