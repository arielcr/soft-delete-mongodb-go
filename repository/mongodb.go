package repository

import (
	"context"
	"errors"
	"time"

	"github.com/arielcr/soft-delete-mongodb-go/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDbRepository mongodb repository struct
type MongoDbRepository struct {
	client *mongo.Client
	config entities.MongoDBConfig
}

// NewMongoDbRepository creates a new mongodb repository
func NewMongoDbRepository(client *mongo.Client, config entities.MongoDBConfig) *MongoDbRepository {
	return &MongoDbRepository{
		client: client,
		config: config,
	}
}

func (r *MongoDbRepository) collection() *mongo.Collection {
	return r.client.Database(r.config.Database).Collection(r.config.Collection)
}

// CreateUser creates a user on the database
func (r *MongoDbRepository) CreateUser(ctx context.Context, user entities.User) (primitive.ObjectID, error) {
	user.ID = primitive.NewObjectID()

	_, err := r.collection().InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return user.ID, nil
}

// GetUser get the user from the database
func (r *MongoDbRepository) GetUser(ctx context.Context, userID primitive.ObjectID) (entities.User, error) {
	var user entities.User

	filter := bson.M{
		"_id": userID,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}

	err := r.collection().FindOne(ctx, filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return user, errors.New("record does not exist")
	} else if err != nil {
		return user, err
	}

	return user, nil
}

// DeleteUser deletes the user from the database
func (r *MongoDbRepository) DeleteUser(ctx context.Context, userID primitive.ObjectID) error {
	filter := bson.M{
		"_id": userID,
		"deleted_at": bson.M{
			"$exists": false,
		},
	}
	updater := bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{
					Key:   "deleted_at",
					Value: time.Now(),
				},
			},
		},
	}

	result, err := r.collection().UpdateOne(ctx, filter, updater)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("record does not exist")
	}

	return nil
}
