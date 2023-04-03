package task

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Repository interface {
	CreateTask(ctx context.Context, task *NewTask) (string, error)
	GetTaskByID(ctx context.Context, id string) (*Task, error)
	UpdateTask(ctx context.Context, task *Task) (*Task, error)
	DeleteTaskByID(ctx context.Context, id string) error
}

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

func (r *repository) CreateTask(ctx context.Context, task *NewTask) (string, error) {
	uOID, err := primitive.ObjectIDFromHex(task.UserID)
	if err != nil {
		return "", err
	}

	doc := bson.M{
		"name":      task.Name,
		"userId":    uOID,
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

func (r *repository) GetTaskByID(ctx context.Context, id string) (*Task, error) {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": oID}
	var taskDocument TaskDocument
	if err := r.coll.FindOne(ctx, &filter).Decode(&taskDocument); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	return &Task{
		ID:        taskDocument.ID.Hex(),
		UserID:    taskDocument.UserID.Hex(),
		Name:      taskDocument.Name,
		Details:   taskDocument.Details,
		Priority:  taskDocument.Priority,
		Category:  taskDocument.Category,
		CreatedAt: taskDocument.CreatedAt,
		UpdatedAt: taskDocument.UpdatedAt,
	}, nil
}

func (r *repository) UpdateTask(ctx context.Context, task *Task) (*Task, error) {
	oID, err := primitive.ObjectIDFromHex(task.ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": oID}
	update := bson.M{
		"$set": &TaskDocument{
			Name:      task.Name,
			Details:   task.Details,
			Priority:  task.Priority,
			Category:  task.Category,
			UpdatedAt: task.UpdatedAt,
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	result := r.coll.FindOneAndUpdate(ctx, &filter, &update, opts)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	var td TaskDocument
	if err := result.Decode(&td); err != nil {
		return nil, err
	}

	return &Task{
		ID:        td.ID.Hex(),
		UserID:    td.UserID.Hex(),
		Name:      td.Name,
		Details:   td.Details,
		Priority:  td.Priority,
		Category:  td.Category,
		CreatedAt: td.CreatedAt,
		UpdatedAt: td.UpdatedAt,
	}, nil
}

func (r *repository) DeleteTaskByID(ctx context.Context, id string) error {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.coll.DeleteOne(ctx, bson.M{"_id": oID})
	if err != nil {
		return err
	}

	// This will be 0 if the document is not found (a 'mongo.ErrNoDocuments' error is not returned)
	deletedCount := result.DeletedCount
	if deletedCount == 0 {
		return ErrTaskNotFound
	}

	return nil
}
