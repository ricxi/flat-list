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
	createTask(ctx context.Context, task *NewTask) (string, error)
	getTaskByID(ctx context.Context, id string) (*Task, error)
	updateTask(ctx context.Context, task *Task) (*Task, error)
	deleteTaskByID(ctx context.Context, id string) error
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

func (r *repository) createTask(ctx context.Context, newTask *NewTask) (string, error) {
	uOID, err := primitive.ObjectIDFromHex(newTask.UserID)
	if err != nil {
		return "", err
	}

	newTaskDoc := NewTaskDocument{
		UserID:    uOID,
		Name:      newTask.Name,
		Details:   newTask.Details,
		Priority:  newTask.Priority,
		Category:  newTask.Category,
		CreatedAt: newTask.CreatedAt,
		UpdatedAt: newTask.UpdatedAt,
	}
	result, err := r.coll.InsertOne(ctx, &newTaskDoc)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *repository) getTaskByID(ctx context.Context, id string) (*Task, error) {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": oID}
	var taskDoc TaskDocument
	if err := r.coll.FindOne(ctx, &filter).Decode(&taskDoc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	return &Task{
		ID:        taskDoc.ID.Hex(),
		UserID:    taskDoc.UserID.Hex(),
		Name:      taskDoc.Name,
		Details:   taskDoc.Details,
		Priority:  taskDoc.Priority,
		Category:  taskDoc.Category,
		CreatedAt: taskDoc.CreatedAt,
		UpdatedAt: taskDoc.UpdatedAt,
	}, nil
}

func (r *repository) updateTask(ctx context.Context, task *Task) (*Task, error) {
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

	var taskDoc TaskDocument
	if err := result.Decode(&taskDoc); err != nil {
		return nil, err
	}

	return &Task{
		ID:        taskDoc.ID.Hex(),
		UserID:    taskDoc.UserID.Hex(),
		Name:      taskDoc.Name,
		Details:   taskDoc.Details,
		Priority:  taskDoc.Priority,
		Category:  taskDoc.Category,
		CreatedAt: taskDoc.CreatedAt,
		UpdatedAt: taskDoc.UpdatedAt,
	}, nil
}

func (r *repository) deleteTaskByID(ctx context.Context, id string) error {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.coll.DeleteOne(ctx, bson.M{"_id": oID})
	if err != nil {
		return err
	}

	// This will be 0 if the document is not found (the mongo driver does not return a 'mongo.ErrNoDocuments' error)
	deletedCount := result.DeletedCount
	if deletedCount == 0 {
		return ErrTaskNotFound
	}

	return nil
}
