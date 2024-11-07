package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Collection

func init() {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	mongoUri := os.Getenv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db = client.Database("backend-gin").Collection("mahasiswa")
	if db == nil {
		log.Fatal("MongoDB collection 'mahasiswa' is not initialized")
	}
}

func GetDb() *mongo.Collection {
	if db == nil {
		log.Fatal("Database connection is not initialized. Call InitDB() first.")
	}
	return db
}