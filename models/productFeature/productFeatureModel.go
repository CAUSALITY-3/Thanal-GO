package productFeatureModel

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductFeature struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	Family    string                 `bson:"family" json:"family" validate:"required,unique"`
	Features  map[string]interface{} `bson:"features" json:"features"`
	CreatedAt time.Time              `bson:"createdAt" json:"createdAt" immutable:"true"`
	UpdatedAt time.Time              `bson:"updatedAt" json:"updatedAt"`
}

func setupProductFeatureSchemaAndIndexes(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define the JSON schema for validation
	schema := bson.M{
		"bsonType": "object",
		"required": []string{"family", "features", "createdAt"},
		"properties": bson.M{
			"id":        bson.M{"bsonType": "objectId"},
			"family":    bson.M{"bsonType": "string"},
			"features":  bson.M{"bsonType": "object"},
			"createdAt": bson.M{"bsonType": "date"},
			"updatedAt": bson.M{"bsonType": "date"},
		},
	}

	validator := bson.M{
		"$jsonSchema": schema,
	}

	// Apply the schema validation
	err := client.Database("test").RunCommand(ctx, bson.D{
		{"collMod", "productFeatures"},
		{"validator", validator},
		{"validationLevel", "strict"},
	}).Err()
	if err != nil {
		log.Printf("Error applying schema validation: %v", err)
	}

	// Create unique indexes for the Family field
	collection := client.Database("test").Collection("productFeatures")
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"family": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err = collection.Indexes().CreateMany(ctx, indexes)
	return err
}
