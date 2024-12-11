package db

import (
	"context"
	"fmt"
	"rest-api/internal/user"
	"rest-api/pkg/logging"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type db struct {
	collection *mongo.Collection
	logger     *logging.Logger
}

// Create implements user.Storage.
func (d *db) Create(ctx context.Context, user user.User) (string, error) {
	//nCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	d.logger.Debug("create user")
	result, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to error: %v", err)
	}
	d.logger.Debug("convert insertedID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	d.logger.Trace(user)
	return "", fmt.Errorf("filed to convert objectid to hex")
}

// Delete implements user.Storage.
func (d *db) Delete(ctx context.Context, id string) error {
	panic("unimplemented")
}

// FindOne implements user.Storage.
func (d *db) FindOne(ctx context.Context, id string) (user.User, error) {
	panic("unimplemented")
}

// Update implements user.Storage.
func (d *db) Update(ctx context.Context, user user.User) error {
	panic("unimplemented")
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {

	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
