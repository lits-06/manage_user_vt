package models

type Token struct {
	ID		string
	Email	string	`bson:"email"`
	Token	string	`bson:"token"`
}