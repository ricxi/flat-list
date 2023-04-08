package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserDocument is used to be
// used by a mongo repository
// implementation to transfer user data
type UserDocument struct {
	OID            primitive.ObjectID `bson:"_id,omitempty"`
	FirstName      string             `bson:"firstName"`
	LastName       string             `bson:"lastName"`
	Email          string             `bson:"email"`
	HashedPassword string             `bson:"hashedPassword"`
	Activated      bool               `bson:"activated"`
	CreatedAt      *time.Time         `bson:"createdAt"`
	UpdatedAt      *time.Time         `bson:"updatedAt"`
}
