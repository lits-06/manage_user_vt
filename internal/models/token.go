package models

import "github.com/dgrijalva/jwt-go"

type Token struct {
	ID		string
	Email	string	`bson:"email"`
	Token	string	`bson:"token"`
}

type Claim struct {
	jwt.StandardClaims
	Email		string	`json:"email"`
}