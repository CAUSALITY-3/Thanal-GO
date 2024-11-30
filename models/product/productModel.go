package productModel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
