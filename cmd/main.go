package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lits-06/manage-user/internal/db/mongodb"
	"github.com/lits-06/manage-user/internal/db/scylla"
	"github.com/lits-06/manage-user/internal/routes"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	scylla.ConnectDB()
	defer scylla.Session.Close()
	mongoClient := mongodb.ConnectDB()
	defer mongoClient.Disconnect(context.Background())

	r := gin.Default()
    routes.SetupRoutes(r)
	r.Run()
}