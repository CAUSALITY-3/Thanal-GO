package userModel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Name            string             `bson:"name" json:"name" validate:"required"`
	Email           string             `bson:"email" json:"email" validate:"required,email" index:"unique"`
	Phone           string             `bson:"phone" json:"phone" index:"unique"`
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
