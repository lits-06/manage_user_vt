package mongodb

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Database *mongo.Database

func ConnectDB() {
	uri := os.Getenv("MONGO_URI")
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Printf("error in connect mongo: %v", err)
	}
	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
        log.Printf("error in ping mongo: %v", err)
    }

	Database = client.Database(os.Getenv("MONGO_DB_NAME"))
}