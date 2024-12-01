package productModel

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Feature struct {
	Type  string `bson:"type" json:"type"`
	Value string `bson:"value" json:"value"`
}

type Rating struct {
	Average float64 `bson:"average" json:"average" default:"0"`
	Count   int     `bson:"count" json:"count" default:"0"`
}

type Review struct {
	Customer   Customer  `bson:"customer" json:"customer"`
	Rating     float64   `bson:"rating" json:"rating"`
	Review     string    `bson:"review" json:"review"`
	ReviewDate time.Time `bson:"reviewDate" json:"reviewDate"`
}

type Customer struct {
	Name  string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
}

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Category    string             `bson:"category" json:"category" validate:"required"`
	Priority    int                `bson:"priority" json:"priority" default:"1"`
	Name        string             `bson:"name" json:"name" validate:"required" index:"unique"`
	Family      string             `bson:"family" json:"family" validate:"required"`
	Description string             `bson:"description" json:"description" validate:"required"`
	Image       string             `bson:"image" json:"image" validate:"required"`
	Images      []string           `bson:"images" json:"images"`
	Price       float64            `bson:"price" json:"price" validate:"required,min=0"`
	Stock       int                `bson:"stock" json:"stock" validate:"required,min=0"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt" immutable:"true"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	VideoURL    string             `bson:"videoUrl,omitempty" json:"videoUrl,omitempty"`
	Features    []Feature          `bson:"features" json:"features"`
	Ratings     Rating             `bson:"ratings" json:"ratings"`
	Reviews     []Review           `bson:"reviews" json:"reviews"`
}

func setupProductSchemaAndIndexes(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define the JSON schema for validation
	schema := bson.M{
		"bsonType": "object",
		"required": []string{"category", "name", "family", "description", "price", "stock", "createdAt"},
		"properties": bson.M{
			"_id":         bson.M{"bsonType": "objectId"},
			"category":    bson.M{"bsonType": "string"},
			"priority":    bson.M{"bsonType": "int", "minimum": 1},
			"name":        bson.M{"bsonType": "string"},
			"family":      bson.M{"bsonType": "string"},
			"description": bson.M{"bsonType": "string"},
			"image":       bson.M{"bsonType": "string"},
			"images": bson.M{
				"bsonType": "array",
				"items":    bson.M{"bsonType": "string"},
			},
			"price": bson.M{
				"bsonType": "double",
				"minimum":  0,
			},
			"stock": bson.M{
				"bsonType": "int",
				"minimum":  0,
			},
			"createdAt": bson.M{"bsonType": "date"},
			"updatedAt": bson.M{"bsonType": "date"},
			"videoUrl":  bson.M{"bsonType": "string"},
			"features": bson.M{
				"bsonType": "array",
				"items": bson.M{
					"bsonType": "object",
					"properties": bson.M{
						"type":  bson.M{"bsonType": "string"},
						"value": bson.M{"bsonType": "string"},
					},
				},
			},
			"ratings": bson.M{
				"bsonType": "object",
				"properties": bson.M{
					"average": bson.M{"bsonType": "double"},
					"count":   bson.M{"bsonType": "int"},
				},
			},
			"reviews": bson.M{
				"bsonType": "array",
				"items": bson.M{
					"bsonType": "object",
					"properties": bson.M{
						"customer": bson.M{
							"bsonType": "object",
							"properties": bson.M{
								"name":  bson.M{"bsonType": "string"},
								"email": bson.M{"bsonType": "string"},
							},
						},
						"rating":     bson.M{"bsonType": "double"},
						"review":     bson.M{"bsonType": "string"},
						"reviewDate": bson.M{"bsonType": "date"},
					},
				},
			},
		},
	}

	validator := bson.M{
		"$jsonSchema": schema,
	}

	// Apply the schema validation
	err := client.Database("test").RunCommand(ctx, bson.D{
		{"collMod", "products"},
		{"validator", validator},
		{"validationLevel", "strict"},
	}).Err()
	if err != nil {
		log.Printf("Error applying schema validation: %v", err)
	}

	// Create unique indexes for the Name field
	collection := client.Database("test").Collection("products")
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"name": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err = collection.Indexes().CreateMany(ctx, indexes)
	return err
}
