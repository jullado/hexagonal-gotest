package repositories

import (
	"context"
	"errors"
	"hexagonal-gotest/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepo struct {
	db         *mongo.Database
	collection string
}

func NewUserRepository(db *mongo.Database, collection string) UserRepository {
	return userRepo{db, collection}
}

func (r userRepo) Gets(filter models.RepoGetUserModel) (result []models.RepoUserModel, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := r.db.Collection(r.collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r userRepo) Create(payload models.RepoCreateUserModel) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = r.db.Collection(r.collection).InsertOne(ctx, payload)
	if err != nil {
		return err
	}

	return nil
}

func (r userRepo) Update(userId string, payload models.RepoUpdateUserModel) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.db.Collection(r.collection).UpdateOne(ctx, bson.D{{Key: "user_id", Value: userId}}, bson.D{{Key: "$set", Value: payload}})
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New(models.ErrUserIdIsNotExist)
	}

	return nil
}

func (r userRepo) Delete(userId string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.db.Collection(r.collection).DeleteOne(ctx, bson.D{{Key: "user_id", Value: userId}})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New(models.ErrUserIdIsNotExist)
	}

	return nil
}
