package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient is the global variable for MongoDB client
var MongoClient *mongo.Client
var Ctx context.Context

// InitMongoDB initializes the MongoDB client
func InitMongoDB() error {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load the connection string from the environment variables
	uri := os.Getenv("MONGO_DB_ATLAS_CREDENTIALS")
	if uri == "" {
		log.Fatal("MONGO_DB_ATLAS_CREDENTIALS not found in environment variables")
	}

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	// Create a new MongoDB client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database to verify the connection
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	// Set the global MongoClient and context
	MongoClient = client
	Ctx = context.Background()

	fmt.Println("Successfully connected to MongoDB!")
	return nil
}
