package task

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Repository interface{}

type repository struct {
	client *mongo.Client
	coll   *mongo.Collection
	dbname string
}

func NewRepository(client *mongo.Client, dbname string) *repository {
	coll := client.Database(dbname).Collection("tasks")

	return &repository{
		client: client,
		coll:   coll,
		dbname: dbname,
	}
}

func NewMongoClient(uri string, timeout int64) (*mongo.Client, error) {
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

func (r *repository) CreateOne(ctx context.Context, task *NewTask) (string, error) {
	uOID, err := primitive.ObjectIDFromHex(task.UserID)
	if err != nil {
		return "", err
	}

	doc := bson.M{
		"name":      task.Name,
		"userID":    uOID,
		"details":   task.Details,
		"priority":  task.Priority,
		"category":  task.Category,
		"createdAt": task.CreatedAt,
		"updatedAt": task.UpdatedAt,
	}
	result, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
