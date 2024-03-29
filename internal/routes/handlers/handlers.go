package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/lits-06/manage-user/internal/db/mongodb"
	"github.com/lits-06/manage-user/internal/db/scylla"
	"github.com/lits-06/manage-user/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	filter := bson.M{"email": user.Email}
	if err := mongodb.UserCollection.FindOne(context.Background(), filter).Err(); err != mongo.ErrNoDocuments {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email has been registed"})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	user.Password = string(hashPassword)

	_, err = mongodb.UserCollection.InsertOne(context.Background(), user)
	if err != nil {
        log.Fatal(err)
    }
	
	c.JSON(http.StatusCreated, gin.H{"message": "Register success"})
}

func Login(c *gin.Context) {
	isNewSession, _ := c.Get("isNewSession")
	if !isNewSession.(bool) {
		c.JSON(http.StatusOK, gin.H{"message": "you have logged in"})
		return
	}

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	var existUser models.User
	filter := bson.M{"email": user.Email}
	if err := mongodb.UserCollection.FindOne(context.Background(), filter).Decode(&existUser); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email has not been registered"})
        return
	}

	err := bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password"})
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Minute).Unix()

	secretKey := []byte(os.Getenv("SECRETKEY"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
        return
    }

	newToken := models.Token{
		Email: user.Email,
		Token: tokenString,
	}

	_, err = mongodb.TokenCollection.InsertOne(context.Background(), newToken)
	if err != nil {
		log.Printf("%v", err)
	}

	c.Header("Authorization", fmt.Sprintf("Bearer %s", tokenString))

	id := gocql.TimeUUID()
	info := fmt.Sprintf("%s login at %s", user.Email, time.Now().Format("02-01-2006 15:04:05"))
	err = scylla.Session.Query(`
		INSERT INTO info (id, email, info)
		VALUES (?, ?, ?)`,
		id, user.Email, info).Exec()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Logout(c *gin.Context) {
	isNewSession, _ := c.Get("isNewSession")
	if isNewSession.(bool) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "you did not log in"})
		return
	}

	userEmail, _ := c.Get("userEmail")
	c.Set("Authorization", "")

	filter := bson.M{"email": userEmail}
    if _, err := mongodb.TokenCollection.DeleteMany(context.Background(), filter); err != nil {
        log.Printf("%v", err)
    }

	id := gocql.TimeUUID()
	info := fmt.Sprintf("%s logout at %s", userEmail, time.Now().Format("02-01-2006 15:04:05"))
	err := scylla.Session.Query(`
		INSERT INTO info (id, email, info)
		VALUES (?, ?, ?)`,
		id, userEmail, info).Exec()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout success"})
}

func Showinfo(c *gin.Context) {
	isNewSession, _ := c.Get("isNewSession")
	if isNewSession.(bool) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "you did not log in"})
		return
	}

	userEmail, _ := c.Get("userEmail")

	var info string
	var infoRecord []string
	query := "SELECT info FROM info WHERE email = ? ALLOW FILTERING"
	iter := scylla.Session.Query(query, userEmail).Iter()
	defer iter.Close()

	for iter.Scan(&info) {
		infoRecord = append(infoRecord, info)
	}

	if err := iter.Close(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
    }

	c.JSON(http.StatusOK, gin.H{"info": infoRecord})
}