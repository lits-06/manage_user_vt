package models

type User struct {
	ID			string
	Email		string	`json:"email" bson:"email"`
	Name		string	`json:"username" bson:"username"`
	Password	string	`json:"password" bson:"password"`
}