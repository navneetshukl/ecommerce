package database

import (
	"context"
	"ecommerce/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	Name  string
	Email string
	Age   int
}

// var Client *mongo.Client

func (DB *Mongo) ConnectDB() *mongo.Client {
	err := godotenv.Load()
	if err != nil {
		DB.client = nil
		log.Fatal("Error loading .env file")
		return nil
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		DB.client = nil
		log.Fatal("MONGO_URI environment variable is not set")
		return nil
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		DB.client = nil
		log.Fatal("Error creating MongoDB client:", err)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to the MongoDB server
	err = client.Connect(ctx)
	if err != nil {
		DB.client = nil
		log.Fatal("Error connecting to MongoDB:", err)
		return nil
	}

	// Ping the MongoDB server to check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		DB.client = nil
		log.Fatal("Error pinging MongoDB:", err)
		return nil
	}

	fmt.Println("Connected to MongoDB")
	DB.client = client
	return client
}

func (DB *Mongo) Insert() {
	database := DB.client.Database("shopping")
	collection := database.Collection("users")

	person := models.User{
		Name:      "Navneet Shukla",
		Email:     "navneetshukla824@gmail.com",
		Password:  "123456",
		Timestamp: time.Now().UTC(),
	}

	insertResult, err := collection.InsertOne(context.Background(), person)
	if err != nil {
		log.Fatal("Error inserting document:", err)
	}

	fmt.Println("Inserted document ID:", insertResult.InsertedID)

}

// CheckUser function return the user based on email of the user
func (DB *Mongo) CheckUser(email string) (bool, error, models.User) {
	database := DB.client.Database("shopping")
	collection := database.Collection("users")

	var user models.User
	fmt.Println("Email from database is ", email)

	err := collection.FindOne(context.Background(), models.User{Email: email}).Decode(&user)

	fmt.Println("USer from databaase is ", user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil, user
		}
		return false, err, user
	}

	return true, nil, user
}

// InsertIntoUser function help in inserting to the database
func (DB *Mongo) InsertIntoUser(user *models.User) (bool, error) {

	database := DB.client.Database("shopping")
	collection := database.Collection("users")
	ID, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal("Error inserting document:", err)
		return false, err
	}

	fmt.Println("ID is ", ID.InsertedID)

	user.ID = ID.InsertedID.(primitive.ObjectID)

	return true, nil

}
