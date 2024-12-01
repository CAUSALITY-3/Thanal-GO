package orderModel

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Review struct {
	Rating     float64   `bson:"rating" json:"rating"`
	Review     string    `bson:"review" json:"review"`
	ReviewDate time.Time `bson:"reviewDate" json:"reviewDate"`
}

type OrderItem struct {
	ProductID    primitive.ObjectID `bson:"productId" json:"productId" validate:"required"`
	ProductName  string             `bson:"productName" json:"productName" validate:"required"`
	ProductImage string             `bson:"productImage" json:"productImage" validate:"required"`
	Quantity     int                `bson:"quantity" json:"quantity" default:"1"`
	Price        float64            `bson:"price" json:"price" validate:"required"`
	Status       string             `bson:"status" json:"status" default:"Waiting for seller approval"`
	Review       Review             `bson:"review" json:"review"`
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt" default:"time.Now()"`
}

type DeliveryAddress struct {
	Name      string `bson:"name" json:"name"`
	HouseName string `bson:"houseName" json:"houseName" validate:"required"`
	Landmark  string `bson:"landmark" json:"landmark" validate:"required"`
	City      string `bson:"city" json:"city" validate:"required"`
	State     string `bson:"state" json:"state" default:""`
	Pincode   int    `bson:"pincode" json:"pincode" validate:"required"`
	Phone     string `bson:"phone" json:"phone" validate:"required"`
}

type Order struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserEmail       string             `bson:"userEmail" json:"userEmail" validate:"required"`
	OrderItems      []OrderItem        `bson:"orderItems" json:"orderItems" validate:"required"`
	TotalPrice      float64            `bson:"totalPrice" json:"totalPrice" validate:"required"`
	Status          string             `bson:"status" json:"status" default:"In progress"`
	DeliveryAddress DeliveryAddress    `bson:"deliveryAddress" json:"deliveryAddress" validate:"required"`
	OrderDate       time.Time          `bson:"orderDate" json:"orderDate" default:"time.Now()"`
}

func setupOrderSchemaAndIndexes(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Define the JSON schema for validation
	schema := bson.M{
		"bsonType": "object",
		"required": []string{"userEmail", "orderItems", "totalPrice", "deliveryAddress", "orderDate"},
		"properties": bson.M{
			"_id":       bson.M{"bsonType": "objectId"},
			"userEmail": bson.M{"bsonType": "string"},
			"orderItems": bson.M{
				"bsonType": "array",
				"items": bson.M{
					"bsonType": "object",
					"properties": bson.M{
						"productId":    bson.M{"bsonType": "objectId"},
						"productName":  bson.M{"bsonType": "string"},
						"productImage": bson.M{"bsonType": "string"},
						"quantity":     bson.M{"bsonType": "int", "minimum": 1},
						"price":        bson.M{"bsonType": "double"},
						"status":       bson.M{"bsonType": "string"},
						"review": bson.M{
							"bsonType": "object",
							"properties": bson.M{
								"rating":     bson.M{"bsonType": "double"},
								"review":     bson.M{"bsonType": "string"},
								"reviewDate": bson.M{"bsonType": "date"},
							},
						},
						"updatedAt": bson.M{"bsonType": "date"},
					},
				},
			},
			"totalPrice": bson.M{"bsonType": "double"},
			"status":     bson.M{"bsonType": "string"},
			"deliveryAddress": bson.M{
				"bsonType": "object",
				"properties": bson.M{
					"name":      bson.M{"bsonType": "string"},
					"houseName": bson.M{"bsonType": "string"},
					"landmark":  bson.M{"bsonType": "string"},
					"city":      bson.M{"bsonType": "string"},
					"state":     bson.M{"bsonType": "string"},
					"pincode":   bson.M{"bsonType": "int"},
					"phone":     bson.M{"bsonType": "string"},
				},
			},
			"orderDate": bson.M{"bsonType": "date"},
		},
	}

	validator := bson.M{
		"$jsonSchema": schema,
	}

	// Apply the schema validation
	err := client.Database("test").RunCommand(ctx, bson.D{
		{"collMod", "orders"},
		{"validator", validator},
		{"validationLevel", "strict"},
	}).Err()
	if err != nil {
		log.Printf("Error applying schema validation: %v", err)
	}

	// Create indexes for the Order collection
	collection := client.Database("test").Collection("orders")
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"userEmail": 1},
			Options: options.Index().SetUnique(false), // Not unique but indexed for faster queries
		},
	}

	_, err = collection.Indexes().CreateMany(ctx, indexes)
	return err
}
