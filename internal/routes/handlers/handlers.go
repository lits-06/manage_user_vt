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

var userCollection = mongodb.Database.Collection(os.Getenv("MONGO_USER_COLLECTION"))
var tokenCollection = mongodb.Database.Collection(os.Getenv("MONGO_TOKEN_COLLECTION"))

func Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	filter := bson.M{"email": user.Email}
	if err := userCollection.FindOne(context.Background(), filter).Err(); err != mongo.ErrNoDocuments {
		c.JSON(http.StatusOK, gin.H{"message": "Email has been registed"})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	user.Password = string(hashPassword)

	_, err = userCollection.InsertOne(context.Background(), user)
	if err != nil {
        log.Fatal(err)
    }
	
	c.JSON(http.StatusCreated, gin.H{"message": "Register success"})
}

func Login(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	filter := bson.M{"email": user.Email}
	res := userCollection.FindOne(context.Background(), filter)
	if res.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email has not been registered"})
        return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	if user.Password != string(hashPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect password"})
        return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(os.Getenv("SECRETKEY"))
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
        return
    }

	newToken := models.Token{
		Email: user.Email,
		Token: tokenString,
	}

	_, err = tokenCollection.InsertOne(context.Background(), newToken)
	if err != nil {
		log.Fatal(err)
	}

	id := gocql.TimeUUID()
	info := fmt.Sprintf("%s login at %s", user.Email, time.Now().Format("02-01-2006 15:04:05"))
	err = scylla.Session.Query(`
		INSERT INTO history (id, email, info)
		VALUES (?, ?, ?)`,
		id, user.Email, info).Exec()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login success"})
}

func Logout(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	filter := bson.M{"email": user.Email}
    if _, err := tokenCollection.DeleteMany(context.Background(), filter); err != nil {
        log.Fatal(err)
    }

	id := gocql.TimeUUID()
	info := fmt.Sprintf("%s logout at %s", user.Email, time.Now().Format("02-01-2006 15:04:05"))
	err := scylla.Session.Query(`
		INSERT INTO history (id, email, info)
		VALUES (?, ?, ?)`,
		id, user.Email, info).Exec()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout success"})
}