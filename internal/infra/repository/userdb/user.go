package userdb

import (
	"chatapp/internal/domain/model"
	"chatapp/internal/domain/repository/userrepo"
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

func (r *Repository) Add(_ context.Context, m model.User) error {
	collection := r.db.Collection("users")
	_, err := collection.InsertOne(context.TODO(), m)
	if err != nil {
		return echo.ErrBadRequest
	}
	return nil
}

func (r *Repository) Get(_ context.Context, cmd userrepo.GetCommand) []model.User {
	collection := r.db.Collection("users")
	var users []model.User

	// Create a filter based on cmd
	filter := bson.D{}
	if cmd.ID != nil {
		filter = append(filter, bson.E{Key: "id", Value: *cmd.ID})
	}
	if cmd.Username != nil {
		filter = append(filter, bson.E{Key: "username", Value: *cmd.Username})
	}
	if cmd.Phone != nil {
		filter = append(filter, bson.E{Key: "phone", Value: *cmd.Phone})
	}
	if cmd.Keyword != nil {
		filter = append(filter, bson.E{Key: "username", Value: bson.D{{Key: "$regex", Value: *cmd.Keyword}}})
	}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil
	}
	for cur.Next(context.TODO()) {
		var user model.User
		err := cur.Decode(&user)
		if err != nil {
			return nil
		}
		users = append(users, user)
	}
	return users
}

func (r *Repository) Delete(ctx context.Context, id uint64) error {
	collection := r.db.Collection("users")
	_, err := collection.DeleteOne(context.TODO(), bson.D{{Key: "id", Value: id}})
	if err != nil {
		return echo.ErrBadRequest
	}
	return nil
}

func (r *Repository) Update(_ context.Context, m model.User) error {
	collection := r.db.Collection("users")
	_, err := collection.UpdateOne(context.TODO(), bson.D{{Key: "id", Value: m.ID}}, bson.D{{Key: "$set", Value: m}})
	if err != nil {
		return echo.ErrNotFound
	}
	return nil
}
