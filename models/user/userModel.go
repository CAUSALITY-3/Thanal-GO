package userModel

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Address struct {
	HouseName string `bson:"houseName" json:"houseName" validate:"required"`
	Landmark  string `bson:"landmark" json:"landmark" validate:"required"`
	City      string `bson:"city" json:"city" validate:"required"`
	State     string `bson:"state" json:"state"`
	Pincode   int    `bson:"pincode" json:"pincode" validate:"required"`
}

type DeliveryAddress struct {
	Name      string `bson:"name" json:"name"`
	HouseName string `bson:"houseName" json:"houseName" validate:"required"`
	Landmark  string `bson:"landmark" json:"landmark" validate:"required"`
	City      string `bson:"city" json:"city" validate:"required"`
	State     string `bson:"state" json:"state"`
	Pincode   int    `bson:"pincode" json:"pincode" validate:"required"`
	Phone     string `bson:"phone" json:"phone"`
}

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name            string             `bson:"name" json:"name,omitempty" `
	Email           string             `bson:"email" json:"email" validate:"required,email" index:"unique"`
	Phone           string             `bson:"phone" json:"phone"`
	ProfilePic      string             `bson:"profilePic" json:"profilePic"`
	Password        string             `bson:"password" json:"password"`
	Address         Address            `bson:"address" json:"address"`
	DeliveryAddress []DeliveryAddress  `bson:"deliveryAddress" json:"deliveryAddress"`
	ProfilePicture  string             `bson:"profilePicture" json:"profilePicture"`
	Orders          []string           `bson:"orders" json:"orders"`
	Wishlists       []string           `bson:"wishlists" json:"wishlists"`
	Bag             []string           `bson:"bag" json:"bag"`
	CreatedAt       time.Time          `bson:"createdAt" json:"createdAt" immutable:"true"`
	LastLoggedIn    time.Time          `bson:"lastLoggedIn" json:"lastLoggedIn"`
	UpdatedAt       time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func SetupSchemaAndIndexes(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define the JSON schema for validation
	validator := bson.M{
		"$jsonSchema": bson.M{
			"bsonType": "object",
			"required": []string{"name", "email"},
			"properties": bson.M{
				"name": bson.M{"bsonType": "string", "minLength": 3},
				"email": bson.M{
					"bsonType": "string",
					"pattern":  `^.+@.+\..+$`, // Basic email pattern
				},
				"phone":      bson.M{"bsonType": "string"},
				"profilePic": bson.M{"bsonType": "string"},
				"password":   bson.M{"bsonType": "string"},
				"address":    bson.M{"bsonType": "object"},
				"deliveryAddress": bson.M{
					"bsonType": "array",
					"items":    bson.M{"bsonType": "object"},
				},
				"orders":       bson.M{"bsonType": "array"},
				"wishlists":    bson.M{"bsonType": "array"},
				"bag":          bson.M{"bsonType": "array"},
				"createdAt":    bson.M{"bsonType": "date"},
				"lastLoggedIn": bson.M{"bsonType": "date"},
				"updatedAt":    bson.M{"bsonType": "date"},
			},
		},
	}

	// Apply the schema validation
	err := client.Database("test").RunCommand(ctx, bson.D{
		{"collMod", "users"},
		{"validator", validator},
		{"validationLevel", "strict"},
	}).Err()
	if err != nil {
		log.Printf("Error applying schema validation: %v", err)
	}

	// Create unique indexes
	collection := client.Database("test").Collection("users")
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err = collection.Indexes().CreateMany(ctx, indexes)
	return err
}
