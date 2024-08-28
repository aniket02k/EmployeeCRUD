package router

import (
	"employeecurd/handlers"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    
    router.HandleFunc("/assignment/user", handlers.GetUser).Methods("GET")
    router.HandleFunc("/assignment/user", handlers.CreateUser).Methods("POST")
    router.HandleFunc("/assignment/user", handlers.UpdateUser).Methods("PATCH")

    return router
}
