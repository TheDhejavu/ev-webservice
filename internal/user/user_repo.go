package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	customErr "github.com/workspace/evoting/ev-webservice/internal/errors"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoUserRepository struct {
	Collection *mongo.Collection
	logger     log.Logger
}

var (
	ErrInvalidId = errors.New("Invalid User ID")
)

// NewMongoUserRepository creates a new user repository.
func NewMongoUserRepository(db *mongo.Database, logger log.Logger) entity.UserRepository {

	collection := db.Collection("users")
	return &mongoUserRepository{collection, logger}
}

// Fetch returns the users with the specified filter from Mongo.
func (repo *mongoUserRepository) Fetch(ctx context.Context, filter interface{}) (res []entity.User, err error) {
	Users := []entity.User{}
	if filter == nil {
		filter = bson.M{}
	}
	projection := bson.M{"password": 0}
	cursor, err := repo.Collection.Find(
		ctx,
		filter,
		options.Find().SetProjection(projection),
	)
	if err != nil {
		repo.logger.Errorf("Fetch transaction error: %s", err)
		return Users, err
	}

	for cursor.Next(ctx) {
		row := entity.User{}
		cursor.Decode(&row)
		Users = append(Users, row)
	}

	return Users, nil
}

// GetByID gets the user with the specified Id from mongo.
func (repo *mongoUserRepository) GetByID(ctx context.Context, id string) (entity.User, error) {
	User := entity.User{}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return User, errors.New("Invalid User ID")
	}
	err = repo.Collection.FindOne(ctx, bson.M{"_id": _id}).Decode(&User)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return User, customErr.ErrEntityDoesNotExist
		default:
			repo.logger.Errorf("FindOne transaction error: %s", err)
			return User, err
		}
	}
	return User, nil
}

// GetByEmail gets the users with the specified email from mongo.
func (repo *mongoUserRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	User := entity.User{}

	err := repo.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&User)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return User, customErr.ErrEntityDoesNotExist
		default:
			repo.logger.Errorf("GetByEmail transaction error: %s", err)
			return User, err
		}
	}

	return User, nil
}

// GetByUsername gets the user with the specified username from mongo.
func (repo *mongoUserRepository) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	User := entity.User{}

	err := repo.Collection.FindOne(ctx, bson.M{"username": username}).Decode(&User)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return User, customErr.ErrEntityDoesNotExist
		default:
			repo.logger.Errorf("GetByEmail transaction error: %s", err)
			return User, err
		}
	}

	return User, nil
}

func constructNotEqualQuery(data map[string]interface{}, query bson.M) (bson.M, error) {
	for field, value := range data {
		if field == "id" {
			id := fmt.Sprintf("%v", value)
			_id, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return query, ErrInvalidId
			}
			query["_id"] = bson.M{"$ne": _id}
		} else {
			query[field] = bson.M{"$ne": value}
		}
	}

	return query, nil
}
func constructQuery(data map[string]interface{}) (bson.M, error) {
	query := bson.M{}

	for field, value := range data {

		if field == "id" {
			id := fmt.Sprintf("%v", value)
			_id, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return query, ErrInvalidId
			}
			query["_id"] = _id
		} else {
			query[field] = value
		}
	}
	return query, nil
}

// GetWithExclude gets the user with the specified user excluding some other data.
func (repo *mongoUserRepository) GetWithExclude(ctx context.Context, user map[string]interface{}, exclude map[string]interface{}) (entity.User, error) {
	User := entity.User{}
	query, err := constructQuery(user)
	query, err = constructNotEqualQuery(exclude, query)

	if err != nil {
		switch err {
		case ErrInvalidId:
			return User, ErrInvalidId
		default:
			return User, err
		}
	}
	err = repo.Collection.FindOne(ctx, query).Decode(&User)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return User, customErr.ErrEntityDoesNotExist
		default:
			repo.logger.Errorf("GetWithExclude transaction error: %s", err)
			return User, err
		}
	}

	return User, nil
}

// Update updates the user with the specified Id from mongo.
func (repo *mongoUserRepository) Update(ctx context.Context, id string, updateUser map[string]interface{}) (entity.User, error) {
	var result entity.User
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, errors.New("Invalid User ID")
	}

	User, err := repo.GetByID(ctx, id)
	User.UpdatedAt = time.Now()

	if err != nil {
		return result, err
	}
	var existingUser map[string]interface{}
	b, err := json.Marshal(User)
	if err != nil {
		return result, err
	}

	json.Unmarshal(b, &existingUser)

	for k := range updateUser {
		if _, ok := existingUser[k]; ok {
			existingUser[k] = updateUser[k]
		}
	}

	delete(existingUser, "_id")

	_, err = repo.Collection.UpdateOne(ctx, bson.M{"_id": _id}, bson.M{"$set": existingUser})
	if err != nil {
		return result, err
	}

	return repo.GetByID(ctx, id)
}

// Create creates a new user.
func (repo *mongoUserRepository) Create(ctx context.Context, User entity.User) (entity.User, error) {
	newUser := entity.User{}
	User.CreatedAt = time.Now()
	res, err := repo.Collection.InsertOne(ctx, User)
	id := res.InsertedID.(primitive.ObjectID).Hex()

	if err != nil {
		repo.logger.Errorf("Create transaction error: %s", err)
		return newUser, err
	}
	return repo.GetByID(ctx, id)
}

// Delete deletes a user.
func (repo *mongoUserRepository) Delete(ctx context.Context, id string) error {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repo.logger.Errorf("DeleteOne transaction error: %s", err)
		return err
	}
	_, err = repo.Collection.DeleteOne(ctx, bson.M{"_id": _id})

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return customErr.ErrEntityDoesNotExist
		default:
			repo.logger.Errorf("Delete transaction error: %s", err)
			return err
		}
	}

	return err
}
