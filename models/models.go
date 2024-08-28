package models

//User Model
type User struct {
	ID        string `bson:"id"`
	FirstName string `bson:"firstname"`
	LastName  string `bson:"lastname"`
	Email     string `bson:"email"`
}

//Employee Model
type Employee struct {
	ID          string `bson:"id"`
	UserID      string `bson:"userId"`
	Designation string `bson:"designation"`
}
