package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/protobuf/proto"

	"employeecurd/db"
	"employeecurd/employeeproto"
	"employeecurd/models"
	"employeecurd/utils"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
    protoBody := r.URL.Query().Get("proto_body")
    //utils.Logger.Println(r.URL.Query())
    decodedProto, err := base64.StdEncoding.DecodeString(protoBody)
    if err != nil {
        utils.Logger.Println("Invalid proto_body:", err)
        http.Error(w, "Invalid proto_body", http.StatusBadRequest)
        return
    }

    var userRequest employeeproto.UserRequest
    err = proto.Unmarshal(decodedProto, &userRequest)
    //utils.Logger.Println("Decoded Proto is:", decodedProto)
    if err != nil {
        utils.Logger.Println("Failed to unmarshal proto:", err)
        http.Error(w, "Failed to unmarshal proto", http.StatusBadRequest)
        return
    }

    userCollection := db.Client.Database("employee").Collection("UserCollection")
    employeeCollection := db.Client.Database("employee").Collection("EmployeeCollection")

    var user models.User
    err = userCollection.FindOne(context.TODO(), bson.M{"id": userRequest.UserId}).Decode(&user)
    if err != nil {
        utils.Logger.Println("User not found. Please provide correct user id:", err)
        http.Error(w, "User not found. Please provide correct user id", http.StatusNotFound)
        return
    }

    var employee models.Employee
    err = employeeCollection.FindOne(context.TODO(), bson.M{"userId": userRequest.UserId}).Decode(&employee)
    if err != nil {
        utils.Logger.Println("User not found in Employee collection:", err)
        http.Error(w, "User not found in Employee collection", http.StatusNotFound)
        return
    }

    response := employeeproto.UserResponse{
        UserId:      user.ID,
        Firstname:   user.FirstName,
        Lastname:    user.LastName,
        Email:       user.Email,
        EmployeeId:  employee.ID,
        Designation: employee.Designation,
    }

    responseData, err := proto.Marshal(&response)
    if err != nil {
        utils.Logger.Println("Failed to marshal response:", err)
        http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/octet-stream")
    w.Write(responseData)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
    var req employeeproto.CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.Logger.Println("Failed to decode request body:", err)
        http.Error(w, "Failed to decode request body", http.StatusBadRequest)
        return
    }

    userCollection := db.Client.Database("employee").Collection("UserCollection")
    employeeCollection := db.Client.Database("employee").Collection("EmployeeCollection")

    user := models.User{
        ID:        generateID(),
        FirstName: req.Firstname,
        LastName:  req.Lastname,
        Email:     req.Email,
    }

    _, err := userCollection.InsertOne(context.TODO(), user)
    if err != nil {
        utils.Logger.Println("Failed to insert user:", err)
        http.Error(w, "Failed to insert user", http.StatusInternalServerError)
        return
    }

    employee := models.Employee{
        ID:          generateID(),
        UserID:      user.ID,
        Designation: req.Designation,
    }

    _, err = employeeCollection.InsertOne(context.TODO(), employee)
    if err != nil {
        utils.Logger.Println("Failed to insert employee:", err)
        http.Error(w, "Failed to insert employee", http.StatusInternalServerError)
        return
    }

    response := employeeproto.CreateUserResponse{UserId: user.ID}
    responseData, err := proto.Marshal(&response)
    if err != nil {
        utils.Logger.Println("Failed to marshal response:", err)
        http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/octet-stream")
    w.Write(responseData)
}

// func UpdateUser(w http.ResponseWriter, r *http.Request) {
//     var req employeeproto.UpdateUserRequest
//     utils.Logger.Println("Inside UpdateUserFunction")
//     if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
//         utils.Logger.Println("Failed to decode request body:", err)
//         http.Error(w, "Failed to decode request body", http.StatusBadRequest)
//         return
//     }

//     userCollection := db.Client.Database("employee").Collection("UserCollection")

//     filter := bson.M{"id": req.UserId}
//     update := bson.M{"$set": bson.M{"email": req.Email}}

//     _, err := userCollection.UpdateOne(context.TODO(), filter, update)
//     if err != nil {
//         utils.Logger.Println("Failed to update user:", err)
//         http.Error(w, "Failed to update user", http.StatusInternalServerError)
//         return
//     }

//     w.WriteHeader(http.StatusOK)
// }

func UpdateUser(w http.ResponseWriter, r *http.Request) {
    var req employeeproto.UpdateUserRequest

    // Added for debug purpose
    utils.Logger.Println("Inside UpdateUserFunction")

    // decode the request body
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        utils.Logger.Println("Failed to decode request body:", err)
        http.Error(w, "Failed to decode request body", http.StatusBadRequest)
        return
    }

    // validate thee request parameters
    if req.UserId == "" {
        utils.Logger.Println("Invalid request: missing UserId")
        http.Error(w, "UserId is required", http.StatusBadRequest)
        return
    }

    if req.Email == "" {
        utils.Logger.Println("Invalid request: missing Email")
        http.Error(w, "Email is required", http.StatusBadRequest)
        return
    }

    // connect to the db
    userCollection := db.Client.Database("employee").Collection("UserCollection")

    // check if the user exists
    var existingUser bson.M
    err := userCollection.FindOne(context.TODO(), bson.M{"id": req.UserId}).Decode(&existingUser)
    if err == mongo.ErrNoDocuments {
        utils.Logger.Println("User not found:", req.UserId)
        http.Error(w, "User not found", http.StatusNotFound)
        return
    } else if err != nil {
        utils.Logger.Println("Error querying database:", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // define  filter and update
    filter := bson.M{"id": req.UserId}
    update := bson.M{"$set": bson.M{"email": req.Email}}

    // update operation
    _, err = userCollection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        utils.Logger.Println("Failed to update user:", err)
        http.Error(w, "Failed to update user", http.StatusInternalServerError)
        return
    }

    utils.Logger.Println("User details updated successfully")
    w.WriteHeader(http.StatusOK)
}

//genereates random id for user.
func generateID() string {
    return uuid.New().String()
}
