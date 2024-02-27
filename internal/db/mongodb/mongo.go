package mongodb

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection
var TokenCollection *mongo.Collection

func ConnectDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Printf("error in connect mongo: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
        log.Printf("error in ping mongo: %v", err)
    }

	db := client.Database(os.Getenv("MONGO_DB_NAME"))
	UserCollection = db.Collection(os.Getenv("MONGO_USER_COLLECTION"))
	TokenCollection = db.Collection(os.Getenv("MONGO_TOKEN_COLLECTION"))

	return client
}