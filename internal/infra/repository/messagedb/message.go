package messagedb

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/domain/repository/messagerepo"
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

func (r *Repository) Add(_ context.Context, m model.Message) error {
	collection := r.db.Collection("messages")
	_, err := collection.InsertOne(context.TODO(), m)
	if err != nil {
		return echo.ErrBadRequest
	}
	return nil
}

func (r *Repository) Get(_ context.Context, cmd messagerepo.GetCommand) []model.Message {
	collection := r.db.Collection("messages")
	var messages []model.Message

	// not allowed to search just by keyword and content type
	if cmd.ID == nil && cmd.ChatID == nil && cmd.Sender == nil {
		return nil
	}

	// Create a filter based on cmd
	filter := bson.D{}
	if cmd.ID != nil {
		filter = append(filter, bson.E{Key: "id", Value: *cmd.ID})
	}
	if cmd.ChatID != nil {
		filter = append(filter, bson.E{Key: "chatid", Value: *cmd.ChatID})
	}
	if cmd.Sender != nil {
		filter = append(filter, bson.E{Key: "sender", Value: *cmd.Sender})
	}
	if cmd.Keyword != nil {
		filter = append(filter, bson.E{Key: "content", Value: bson.D{{Key: "$regex", Value: *cmd.Keyword}}})
	}
	if cmd.ContentType != nil {
		filter = append(filter, bson.E{Key: "contenttype", Value: *cmd.ContentType})
	}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil
	}
	for cur.Next(context.TODO()) {
		var message model.Message
		err := cur.Decode(&message)
		if err != nil {
			return nil
		}
		messages = append(messages, message)
	}
	return messages
}

func (r *Repository) Delete(_ context.Context, id uint64) error {
	collection := r.db.Collection("messages")
	_, err := collection.DeleteOne(context.TODO(), bson.D{{Key: "id", Value: id}})
	if err != nil {
		return echo.ErrBadRequest
	}
	return nil
}
