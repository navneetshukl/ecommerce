package database

import (
	"ecommerce/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type Mongo struct {
	client *mongo.Client
}

type dbHandler interface {
	ConnectDB() *mongo.Client
	Insert()
	CheckUser(email string) (bool, error)
	InsertIntoUser(user models.User) (bool, error)
}
