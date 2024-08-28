package main

import (
	"log"
	"net/http"

	"employeecurd/db"
	"employeecurd/router"
	"employeecurd/utils"
)

func main() {
    // Initialize the logger
    utils.InitLogger()

    utils.Logger.Println("This is a test log")
    // Connect to MongoDB
    err := db.ConnectToDB("mongodb://localhost:27017")
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }

    // Initialize the router
    r := router.NewRouter()

    // Start the server
    log.Println("Server starting on port 8000")
    if err := http.ListenAndServe(":8000", r); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}
