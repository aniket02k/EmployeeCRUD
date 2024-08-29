package handlers

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"employeecurd/db"
	"employeecurd/employeeproto"
	"employeecurd/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/protobuf/proto"
)

func TestGetUser(t *testing.T) {
    // Setup in-memory MongoDB client
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        t.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    db.Client = client

    // Setup mock data
    userCollection := client.Database("employee").Collection("UserCollection")
    employeeCollection := client.Database("employee").Collection("EmployeeCollection")

    user := models.User{
        ID:        "user123",
        FirstName: "John",
        LastName:  "Doe",
        Email:     "john.doe@example.com",
    }
    _, err = userCollection.InsertOne(context.Background(), user)
    if err != nil {
        t.Fatalf("Failed to insert user: %v", err)
    }

    employee := models.Employee{
        ID:          "emp123",
        UserID:      "user123",
        Designation: "Engineer",
    }
    _, err = employeeCollection.InsertOne(context.Background(), employee)
    if err != nil {
        t.Fatalf("Failed to insert employee: %v", err)
    }

    userRequest := &employeeproto.UserRequest{UserId: "user123"}
    protoData, err := proto.Marshal(userRequest)
    if err != nil {
        t.Fatalf("Failed to marshal proto: %v", err)
    }
    protoBody := base64.StdEncoding.EncodeToString(protoData)

    req, err := http.NewRequest("GET", "/assignment/user?proto_body="+protoBody, nil)
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }

    rr := httptest.NewRecorder()

    // Use the userHandler struct
    handler := NewUserHandler(userCollection, employeeCollection)
    handler.GetUser(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
}

func TestCreateUser(t *testing.T) {
    client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        t.Fatalf("Failed to connect to MongoDB: %v", err)
    }
    db.Client = client

    userCollection := client.Database("employee").Collection("UserCollection")
    employeeCollection := client.Database("employee").Collection("EmployeeCollection")

    reqBody := `{
        "firstname": "Jane",
        "lastname": "Doe",
        "email": "jane.doe@example.com",
        "designation": "Manager"
    }`
    req, err := http.NewRequest("POST", "/assignment/user", strings.NewReader(reqBody))
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")

    rr := httptest.NewRecorder()

    // Use the userHandler struct
    handler := NewUserHandler(userCollection, employeeCollection)
    handler.CreateUser(rr, req)

    assert.Equal(t, http.StatusOK, rr.Code)
}
