package filedb

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/domain/repository/filerepo"
	"context"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client *mongo.Client
	db     *mongo.Database
}

func New() (*Repository, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}
	db := client.Database("chat_app")
	return &Repository{
		client: client,
		db:     db,
	}, nil
}

func (r *Repository) Add(_ context.Context, m model.File) error {
	collection := r.db.Collection("files")
	_, err := collection.InsertOne(context.TODO(), m)
	if err != nil {
		return echo.ErrBadRequest
	}
	return nil
}

func (r *Repository) Get(_ context.Context, cmd filerepo.GetCommand) []model.File {
	collection := r.db.Collection("files")
	var files []model.File

	// Create a filter based on cmd
	filter := bson.D{}
	if cmd.ID != nil {
		filter = append(filter, bson.E{Key: "id", Value: *cmd.ID})
	}
	if cmd.UserID != nil {
		filter = append(filter, bson.E{Key: "userid", Value: *cmd.UserID})
	}
	if cmd.FileName != nil {
		filter = append(filter, bson.E{Key: "filename", Value: *cmd.FileName})
	}
	if cmd.ContentType != nil {
		filter = append(filter, bson.E{Key: "contenttype", Value: *cmd.ContentType})
	}
	if cmd.ChatID != nil {
		filter = append(filter, bson.E{Key: "chatids", Value: bson.D{{Key: "$in", Value: *cmd.ChatID}}})
	}
	if cmd.Keyword != nil {
		filter = append(filter, bson.E{Key: "filename", Value: bson.D{{Key: "$regex", Value: *cmd.Keyword}}})
	}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil
	}
	for cur.Next(context.TODO()) {
		var file model.File
		err := cur.Decode(&file)
		if err != nil {
			return nil
		}
		files = append(files, file)
	}
	return files
}

func (r *Repository) Delete(ctx context.Context, id uint64) error {
	collection := r.db.Collection("files")
	_, err := collection.DeleteOne(context.TODO(), bson.D{{Key: "id", Value: id}})
	if err != nil {
		return echo.ErrBadRequest
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, m model.File) error {
	collection := r.db.Collection("files")
	_, err := collection.UpdateOne(ctx, bson.M{"id": m.ID}, bson.M{"$set": m})
	if err != nil {
		return filerepo.ErrIDNotFound
	}
	return nil
}
