package db

import (
	"context"
	"errors"
	"fmt"
	"rest-api/internal/apperror"
	"rest-api/internal/user"
	"rest-api/pkg/logging"

	"go.mongodb.org/mongo-driver/bson"
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
	return "", fmt.Errorf("filed to convert objectid to hex. probably oid: %s", oid)
}

func (d *db) FindAll(ctx context.Context) (u []user.User, err error) {

	cursor, err := d.collection.Find(ctx, bson.M{})
	if cursor.Err() != nil {
		return u, fmt.Errorf("failed to find all users due to error %v", err)
	}

	if err = cursor.All(ctx, &u); err != nil {
		return u, fmt.Errorf("failed to read all documents from cursor. error %v", err)
	}
	return u, nil
}

// FindOne implements user.Storage.
func (d *db) FindOne(ctx context.Context, id string) (u user.User, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return u, fmt.Errorf("failed to convert hex to objectid. hex: %s", id)
	}
	filter := bson.M{"_id": oid}

	result := d.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, apperror.ErrNotFound
		}
		return u, fmt.Errorf("failed to find one user by id: %s due to error %v", id, err)
	}
	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("failed to decode user(id:%s) from DB due to error %v", id, err)
	}
	return u, nil
}

// Update implements user.Storage.
func (d *db) Update(ctx context.Context, user user.User) error {
	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID. ID=%s", user.ID)
	}
	filter := bson.M{"_id": objectID}
	userBytes, err := bson.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user, err: %v", err)
	}

	var updateUserObj bson.M
	err = bson.Unmarshal(userBytes, &updateUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user, err: %v", err)
	}
	delete(updateUserObj, "_id")

	update := bson.M{
		"$set": updateUserObj,
	}

	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update user query. error: %v", err)
	}
	if result.MatchedCount == 0 {
		return apperror.ErrNotFound
	}

	d.logger.Tracef("Matched %d documents and Modified %d documents", result.MatchedCount, result.ModifiedCount)

	return nil
}

// Delete implements user.Storage.
func (d *db) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert user ID to ObjectID. ID=%s", id)
	}
	filter := bson.M{"_id": objectID}

	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %v", err)
	}
	if result.DeletedCount == 0 {
		return apperror.ErrNotFound
	}

	d.logger.Tracef("Deleted %d documents s", result.DeletedCount)
	return nil
}

func NewStorage(database *mongo.Database, collection string, logger *logging.Logger) user.Storage {

	return &db{
		collection: database.Collection(collection),
		logger:     logger,
	}
}
