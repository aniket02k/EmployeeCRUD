package main

import (
	"log"
	"net/http"

	"employeecurd/db"
	"employeecurd/handlers"
	"employeecurd/router"
	"employeecurd/utils"
)

func main() {
	// Initialize the logger
	utils.InitLogger()

	utils.Logger.Println("Application is starting...")

	// Connect to MongoDB
	err := db.ConnectToDB("mongodb://localhost:27017")
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Get the MongoDB collections
	userCollection := db.Client.Database("employee").Collection("UserCollection")
	employeeCollection := db.Client.Database("employee").Collection("EmployeeCollection")

	// Initialize the userHandler
	userHandler := handlers.NewUserHandler(userCollection, employeeCollection)

	// Initialize the router with the userHandler
	r := router.NewRouter(userHandler)

	// Start the server
	log.Println("Server starting on port 8000")
	if err := http.ListenAndServe(":8000", r); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
