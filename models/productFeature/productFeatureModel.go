package productFeatureModel

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductFeature struct {
	ID        primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	Family    string                 `bson:"family" json:"family" validate:"required,unique"`
	Features  map[string]interface{} `bson:"features" json:"features"`
	CreatedAt time.Time              `bson:"createdAt" json:"createdAt" immutable:"true"`
	UpdatedAt time.Time              `bson:"updatedAt" json:"updatedAt"`
}
