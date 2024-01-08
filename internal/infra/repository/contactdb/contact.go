package contactdb

import (
	"chatapp/internal/domain/model"
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

func (r *Repository) Add(_ context.Context, m model.Contact) error {
	collection := r.db.Collection("contacts")
	_, err := collection.InsertOne(context.TODO(), m)
	if err != nil {
		return echo.ErrBadRequest
	}
	return nil
}

func (r *Repository) Get(_ context.Context, uid uint64) []model.Contact {
	collection := r.db.Collection("contacts")
	var contacts []model.Contact

	// Create a filter based on uid
	filter := bson.D{{Key: "userid", Value: uid}}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil
	}
	for cur.Next(context.TODO()) {
		var contact model.Contact
		err := cur.Decode(&contact)
		if err != nil {
			return nil
		}
		contacts = append(contacts, contact)
	}
	return contacts
}

func (r *Repository) Delete(_ context.Context, uid uint64, cid uint64) error {
	collection := r.db.Collection("contacts")
	_, err := collection.DeleteOne(context.TODO(), bson.D{{Key: "userid", Value: uid}, {Key: "contactid", Value: cid}})
	if err != nil {
		return echo.ErrBadRequest
	}
	return nil
}
