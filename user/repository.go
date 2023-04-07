package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Repository interface {
	CreateUser(ctx context.Context, user UserRegistrationInfo) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*UserInfo, error)
	UpdateUserByID(ctx context.Context, u UserInfo) error
}

// mongoRepository implements Repository interface
type mongoRepository struct {
	client   *mongo.Client
	database string
	coll     *mongo.Collection
	timeout  time.Duration
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
func NewRepository(client *mongo.Client, database string, timeout int) Repository {
	usersCollection := client.Database(database).Collection("users")

	m := mongoRepository{
		client:   client,
		database: database,
		timeout:  time.Duration(timeout) * time.Second,
		coll:     usersCollection,
	}

	m.setupIndexes()

	return &m
}

func (m *mongoRepository) setupIndexes() {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys: bson.D{{
			Key:   "email",
			Value: 1,
		}},
		Options: options.Index().SetUnique(true),
	}

	_, err := m.coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Println("setupIdexes:", err)
	}
}

// CreateOne inserts a new user with a unique email into the database.
func (m *mongoRepository) CreateUser(ctx context.Context, u UserRegistrationInfo) (string, error) {
	userInfo := bson.M{
		"firstName":      u.FirstName,
		"lastName":       u.LastName,
		"email":          u.Email,
		"hashedPassword": u.HashedPassword,
		"activated":      u.Activated,
		"createdAt":      u.CreatedAt,
		"updatedAt":      u.UpdatedAt,
	}
	result, err := m.coll.InsertOne(ctx, &userInfo)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", ErrDuplicateUser
		}
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

// GetUserByEmail Queries a user with their email
func (m *mongoRepository) GetUserByEmail(ctx context.Context, email string) (*UserInfo, error) {
	type UserDoc struct {
		OID            primitive.ObjectID `bson:"_id,omitempty"`
		FirstName      string             `bson:"firstName"`
		LastName       string             `bson:"lastName"`
		Email          string             `bson:"email"`
		HashedPassword string             `bson:"hashedPassword"`
		Activated      bool               `bson:"activated"`
		CreatedAt      *time.Time         `bson:"createdAt"`
		UpdatedAt      *time.Time         `bson:"updatedAt"`
	}

	var userDoc UserDoc
	filter := bson.M{"email": email}
	if err := m.coll.FindOne(ctx, filter).Decode(&userDoc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("%w by email", ErrUserNotFound)
		}
		return nil, err
	}

	return &UserInfo{
		ID:             userDoc.OID.Hex(),
		FirstName:      userDoc.FirstName,
		LastName:       userDoc.LastName,
		Email:          userDoc.Email,
		HashedPassword: userDoc.HashedPassword,
		Activated:      userDoc.Activated,
		CreatedAt:      userDoc.CreatedAt,
		UpdatedAt:      userDoc.UpdatedAt,
	}, nil
}

func (m *mongoRepository) findUserByID(ctx context.Context, id string) (*UserInfo, error) {
	userOID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var userDocument UserDocument
	if err := m.coll.FindOne(ctx, bson.M{"_id": userOID}).Decode(&userDocument); err != nil {
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

// UpdateUserByID updates a user's info based on their id
// ! It's currently only set up to update a user's activation status, but this will change
func (m *mongoRepository) UpdateUserByID(ctx context.Context, u UserInfo) error {
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
	result := m.coll.FindOneAndUpdate(ctx, filter, update)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("unable to update by id: %w", ErrUserNotFound)
		}
		return err
	}

	return nil
}
