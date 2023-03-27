package task

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskDocument is a data type used by
// a repository layer that implements mongo
// to communicate with the mongo API.
// The 'omitempty' tag is used to prevent an empty
// field value from wiping a document during an update.
type TaskDocument struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"userId,omitempty"`
	Name      string             `bson:"name,omitempty"`
	Details   string             `bson:"details,omitempty"`
	Priority  string             `bson:"priority,omitempty"`
	Category  string             `bson:"category,omitempty"`
	CreatedAt *time.Time         `bson:"createdAt,omitempty"`
	UpdatedAt *time.Time         `bson:"updatedAt,omitempty"`
}
