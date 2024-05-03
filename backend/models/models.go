package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name,omitempty" json:"name" binding:"required"`
	Email     string             `bson:"email,omitempty" json:"email" binding:"required"`
	Password  string             `bson:"password,omitempty" json:"password" binding:"required"`
	Phone     string             `bson:"phone,omitempty" json:"phone" binding:"required"`
	Address   string             `bson:"address,omitempty" json:"address" binding:"required"`
	Role      int                `bson:"role,omitempty" json:"role"`
	Timestamp time.Time          `bson:"timestamp,omitempty" json:"timestamp"`
}
