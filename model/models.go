package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:model
type Prime struct {
	// Id of the Movie
	// required: true
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" `

	// Movie Name
	// required: true
	Movie string `json:"movie,omitempty" validate:"required"`

	// Watched check
	// required: true
	Watched bool `json:"watched,omitempty" validate:"required"`
}
