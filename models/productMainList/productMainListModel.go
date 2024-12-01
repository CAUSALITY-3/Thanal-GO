package productMainListsModel

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Rating struct {
	Average float64 `bson:"average" json:"average" validate:"required"`
	Count   int     `bson:"count" json:"count" validate:"required"`
}

type Product struct {
	ProductID   primitive.ObjectID `bson:"productId" json:"productId" validate:"required"`
	Category    string             `bson:"category" json:"category" validate:"required"`
	Name        string             `bson:"name" json:"name" validate:"required"`
	Description string             `bson:"description" json:"description" validate:"required"`
	Price       float64            `bson:"price" json:"price" validate:"required"`
	Image       string             `bson:"image" json:"image" validate:"required"`
	Stock       int                `bson:"stock" json:"stock" validate:"required"`
	Ratings     Rating             `bson:"ratings" json:"ratings" validate:"required"`
}

type ProductMainList struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      string             `bson:"type" json:"type" validate:"required"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt" immutable:"true"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	Data      map[string]Product `bson:"data" json:"data" validate:"required"` // Map of string to Product
}

func setupProductMainListSchemaAndIndexes(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define the JSON schema for validation
	schema := bson.M{
		"bsonType": "object",
		"required": []string{"category", "name", "description", "price", "image", "stock", "ratings"},
		"properties": bson.M{
			"productId":   bson.M{"bsonType": "objectId"},
			"category":    bson.M{"bsonType": "string"},
			"name":        bson.M{"bsonType": "string"},
			"description": bson.M{"bsonType": "string"},
			"price":       bson.M{"bsonType": "double"},
			"image":       bson.M{"bsonType": "string"},
			"stock":       bson.M{"bsonType": "int"},
			"ratings": bson.M{
				"bsonType": "object",
				"properties": bson.M{
					"average": bson.M{"bsonType": "double"},
					"count":   bson.M{"bsonType": "int"},
				},
				"required": []string{"average", "count"},
			},
			"createdAt": bson.M{"bsonType": "date"},
			"updatedAt": bson.M{"bsonType": "date"},
			"data": bson.M{
				"bsonType": "object",
				"additionalProperties": bson.M{
					"bsonType": "object",
				},
			},
		},
	}

	validator := bson.M{
		"$jsonSchema": schema,
	}

	// Apply the schema validation
	err := client.Database("test").RunCommand(ctx, bson.D{
		{"collMod", "productmainlists"},
		{"validator", validator},
		{"validationLevel", "strict"},
	}).Err()
	if err != nil {
		log.Printf("Error applying schema validation: %v", err)
	}

	// Create unique indexes (if needed)
	collection := client.Database("test").Collection("productmainlists")
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"name": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err = collection.Indexes().CreateMany(ctx, indexes)
	return err
}
