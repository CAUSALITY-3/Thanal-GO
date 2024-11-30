package orderModel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
