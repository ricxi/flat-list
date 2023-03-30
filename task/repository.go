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
	DeleteTaskByID(ctx context.Context, id string) (int64, error)
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
	taskDoc := struct {
		ID        string     `bson:"_id,omitempty"`
		UserID    string     `bson:"userId,omitempty"`
		Name      string     `bson:"name"`
		Details   string     `bson:"details,omitempty"`
		Priority  string     `bson:"priority,omitempty"`
		Category  string     `bson:"category,omitempty"`
		CreatedAt *time.Time `bson:"createdAt,omitempty"`
		UpdatedAt *time.Time `bson:"updatedAt,omitempty"`
	}{}
	if err := r.coll.FindOne(ctx, &filter).Decode(&taskDoc); err != nil {
		return nil, err
	}

	// Checking the result for an error, then running the line below causes a nil pointer runtime error if no document is found
	// if err := result.Decode(&task); err != nil {
	// 	return nil, err
	// }

	return &Task{
		ID:        taskDoc.ID,
		UserID:    taskDoc.UserID,
		Name:      taskDoc.Name,
		Details:   taskDoc.Details,
		Priority:  taskDoc.Priority,
		Category:  taskDoc.Category,
		CreatedAt: taskDoc.CreatedAt,
		UpdatedAt: taskDoc.UpdatedAt,
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

func (r *repository) DeleteTaskByID(ctx context.Context, id string) (int64, error) {
	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}

	result, err := r.coll.DeleteOne(ctx, bson.M{"_id": oID})
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}
