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

	"employeecurd/employeeproto"
	"employeecurd/models"
	"employeecurd/utils"
)

type userHandler struct {
	UserCollection     *mongo.Collection
	EmployeeCollection *mongo.Collection
}

func NewUserHandler(userCollection, employeeCollection *mongo.Collection) UserHandler {
	return &userHandler{
		UserCollection:     userCollection,
		EmployeeCollection: employeeCollection,
	}
}

func (h *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	protoBody := r.URL.Query().Get("proto_body")

   // utils.Logger.Println("Inside GetUser handler")
	decodedProto, err := base64.StdEncoding.DecodeString(protoBody)
	if err != nil {
		utils.Logger.Println("Invalid proto_body:", err)
		http.Error(w, "Invalid proto_body", http.StatusBadRequest)
		return
	}

	var userRequest employeeproto.UserRequest
	err = proto.Unmarshal(decodedProto, &userRequest)
	if err != nil {
		utils.Logger.Println("Failed to unmarshal proto:", err)
		http.Error(w, "Failed to unmarshal proto", http.StatusBadRequest)
		return
	}

	var user models.User
	err = h.UserCollection.FindOne(context.TODO(), bson.M{"id": userRequest.UserId}).Decode(&user)
	if err != nil {
		utils.Logger.Println("User not found. Please provide correct user id:", err)
		http.Error(w, "User not found. Please provide correct user id", http.StatusNotFound)
		return
	}

	var employee models.Employee
	err = h.EmployeeCollection.FindOne(context.TODO(), bson.M{"userId": userRequest.UserId}).Decode(&employee)
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

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    //utils.Logger.Println("Inside CreateUser handler")
	var req employeeproto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Logger.Println("Failed to decode request body:", err)
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	user := models.User{
		ID:        generateID(),
		FirstName: req.Firstname,
		LastName:  req.Lastname,
		Email:     req.Email,
	}

	_, err := h.UserCollection.InsertOne(context.TODO(), user)
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

	_, err = h.EmployeeCollection.InsertOne(context.TODO(), employee)
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

func (h *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
    utils.Logger.Println("Inside Update user handler")
	var req employeeproto.UpdateUserRequest

	utils.Logger.Println("Inside UpdateUserFunction")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Logger.Println("Failed to decode request body:", err)
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

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

	var existingUser bson.M
	err := h.UserCollection.FindOne(context.TODO(), bson.M{"id": req.UserId}).Decode(&existingUser)
	if err == mongo.ErrNoDocuments {
		utils.Logger.Println("User not found:", req.UserId)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		utils.Logger.Println("Error querying database:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	filter := bson.M{"id": req.UserId}
	update := bson.M{"$set": bson.M{"email": req.Email}}

	_, err = h.UserCollection.UpdateOne(context.TODO(), filter, update)
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
