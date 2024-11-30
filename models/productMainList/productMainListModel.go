package productModel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
