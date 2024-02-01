package chatdb

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/domain/repository/chatrepo"
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

func (r *Repository) Add(_ context.Context, m model.Chat) error {
	collection := r.db.Collection("chats")
	_, err := collection.InsertOne(context.TODO(), m)
	if err != nil {
		return echo.ErrBadRequest
	}
	return nil
}

func (r *Repository) Get(_ context.Context, cmd chatrepo.GetCommand) []model.Chat {
	collection := r.db.Collection("chats")
	var chats []model.Chat

	// Create a filter based on cmd
	filter := bson.D{}
	if cmd.ID != nil {
		filter = append(filter, bson.E{Key: "id", Value: *cmd.ID})
	}
	if cmd.UserID != nil {
		filter = append(filter, bson.E{Key: "people", Value: bson.D{{Key: "$all", Value: *cmd.UserID}}})
	}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil
	}
	for cur.Next(context.TODO()) {
		var chat model.Chat
		err := cur.Decode(&chat)
		if err != nil {
			return nil
		}
		chats = append(chats, chat)
	}
	return chats
}

func (r *Repository) Delete(ctx context.Context, id uint64) error {
	collection := r.db.Collection("chats")
	_, err := collection.DeleteOne(context.TODO(), bson.D{{Key: "id", Value: id}})
	if err != nil {
		return echo.ErrBadRequest
	}
	return nil
}

func (r *Repository) Update(_ context.Context, m model.Chat) error {
	collection := r.db.Collection("chats")
	_, err := collection.UpdateOne(context.TODO(), bson.D{{Key: "id", Value: m.ID}}, bson.D{{Key: "$set", Value: m}})
	if err != nil {
		return echo.ErrNotFound
	}
	return nil
}
