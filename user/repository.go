package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ! I'm not sure where to convert a UserInfo type to a UserDocument type
// ! in the repository or service layer
type Repository interface {
	createUser(ctx context.Context, user UserRegistrationInfo) (string, error)
	getUserByEmail(ctx context.Context, email string) (*UserInfo, error)
	updateUserByID(ctx context.Context, u UserInfo) error
	getUserByID(ctx context.Context, id string) (*UserInfo, error)
}

// repository implements Repository interface
type repository struct {
	client   *mongo.Client
	database string
	coll     *mongo.Collection
}

func NewMongoClient(uri string, timeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return client, nil
}

// Create a new user repository with the mongo client and database name
func NewRepository(client *mongo.Client, database string) Repository {
	usersCollection := client.Database(database).Collection("users")

	m := repository{
		client:   client,
		database: database,
		coll:     usersCollection,
	}

	return &m
}

// CreateUser inserts a new user with a unique email into the database.
func (r *repository) createUser(ctx context.Context, u UserRegistrationInfo) (string, error) {
	userDocument := UserRegistrationDocument{
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Email:          u.Email,
		HashedPassword: u.HashedPassword,
		Activated:      u.Activated,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
	result, err := r.coll.InsertOne(ctx, &userDocument)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", ErrDuplicateUser
		}
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// getUserByEmail Queries a user with their email
func (r *repository) getUserByEmail(ctx context.Context, email string) (*UserInfo, error) {
	var userDocument UserDocument
	filter := bson.M{"email": email}
	if err := r.coll.FindOne(ctx, filter).Decode(&userDocument); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%w by email", ErrUserNotFound)
		}
		return nil, err
	}

	return &UserInfo{
		ID:             userDocument.OID.Hex(),
		FirstName:      userDocument.FirstName,
		LastName:       userDocument.LastName,
		Email:          userDocument.Email,
		HashedPassword: userDocument.HashedPassword,
		Activated:      userDocument.Activated,
		CreatedAt:      userDocument.CreatedAt,
		UpdatedAt:      userDocument.UpdatedAt,
	}, nil
}

func (r *repository) getUserByID(ctx context.Context, id string) (*UserInfo, error) {
	userOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var userDocument UserDocument
	if err := r.coll.FindOne(ctx, bson.M{"_id": userOID}).Decode(&userDocument); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%w by email", ErrUserNotFound)
		}
		return nil, err
	}

	return &UserInfo{
		ID:             userDocument.OID.Hex(),
		FirstName:      userDocument.FirstName,
		LastName:       userDocument.LastName,
		Email:          userDocument.Email,
		HashedPassword: userDocument.HashedPassword,
		Activated:      userDocument.Activated,
		CreatedAt:      userDocument.CreatedAt,
		UpdatedAt:      userDocument.UpdatedAt,
	}, nil
}

// updateUserByID updates a user's info based on their id
// ! It's currently only set up to update a user's activation status, but this will change
func (r *repository) updateUserByID(ctx context.Context, u UserInfo) error {
	userOID, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": userOID}
	update := bson.M{
		"$set": bson.M{
			"activated": u.Activated,
			"updatedAt": u.UpdatedAt,
		},
	}
	result := r.coll.FindOneAndUpdate(ctx, filter, update)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("unable to update by id: %w", ErrUserNotFound)
		}
		return err
	}

	return nil
}
