package user

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepository struct {
	Collection *mongo.Collection
	logger     log.Logger
}

func NewMongoUserRepository(db *mongo.Database, logger log.Logger) entity.UserRepository {

	collection := db.Collection("Users")
	return &mongoUserRepository{collection, logger}
}

func (repo *mongoUserRepository) Fetch(ctx context.Context, filter interface{}) (res []entity.User, err error) {
	Users := []entity.User{}
	if filter == nil {
		filter = bson.M{}
	}

	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		repo.logger.Errorf("Find transaction error: %s", err)
		return Users, err
	}

	for cursor.Next(ctx) {
		row := entity.User{}
		cursor.Decode(&row)
		Users = append(Users, row)
	}

	return Users, nil
}
func (repo *mongoUserRepository) GetByID(ctx context.Context, id string) (entity.User, error) {
	User := entity.User{}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return User, errors.New("Invalid User ID")
	}
	err = repo.Collection.FindOne(ctx, bson.M{"_id": _id}).Decode(&User)
	if err != nil {
		repo.logger.Errorf("GetByID transaction error: %s", err)
		return User, err
	}
	return User, nil
}
func (repo *mongoUserRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	User := entity.User{}

	err := repo.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&User)

	if err != nil {
		repo.logger.Errorf("GetByEmail transaction error: %s", err)
		return User, err
	}

	return User, nil
}
func (repo *mongoUserRepository) Update(ctx context.Context, id string, data interface{}) (entity.User, error) {
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
	var exist map[string]interface{}
	b, err := json.Marshal(User)
	if err != nil {
		return result, err
	}

	json.Unmarshal(b, &exist)
	var change map[string]interface{}

	b, err = json.Marshal(data)
	if err != nil {
		return result, err
	}
	json.Unmarshal(b, &change)

	for k := range change {
		if _, ok := exist[k]; ok {
			exist[k] = change[k]
		}
	}

	delete(exist, "_id")

	_, err = repo.Collection.UpdateOne(ctx, bson.M{"_id": _id}, bson.M{"$set": exist})
	if err != nil {
		return result, err
	}

	return repo.GetByID(ctx, id)
}
func (repo *mongoUserRepository) Store(ctx context.Context, User entity.User) (entity.User, error) {
	newUser := entity.User{}
	User.CreatedAt = time.Now()
	res, err := repo.Collection.InsertOne(ctx, User)
	id := res.InsertedID.(primitive.ObjectID).Hex()

	if err != nil {
		repo.logger.Errorf("Store transaction error: %s", err)
		return newUser, err
	}
	return repo.GetByID(ctx, id)
}
func (repo *mongoUserRepository) Delete(ctx context.Context, id string) error {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repo.logger.Errorf("DeleteOne transaction error: %s", err)
		return err
	}
	_, err = repo.Collection.DeleteOne(ctx, bson.M{"_id": _id})
	if err != nil {
		repo.logger.Errorf("DeleteOne transaction error: %s", err)
		return nil
	}
	return err
}
