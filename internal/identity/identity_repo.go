package identity

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

type mongoIdentityRepository struct {
	Collection *mongo.Collection
	logger     log.Logger
}

func NewMongoIdentityRepository(db *mongo.Database, logger log.Logger) entity.IdentityRepository {

	collection := db.Collection("Identitys")
	return &mongoIdentityRepository{collection, logger}
}

func (repo *mongoIdentityRepository) Fetch(ctx context.Context, filter interface{}) (res []entity.Identity, err error) {
	Identitys := []entity.Identity{}
	if filter == nil {
		filter = bson.M{}
	}

	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		repo.logger.Errorf("Fetch transaction error: %s", err)
		return Identitys, err
	}

	for cursor.Next(ctx) {
		row := entity.Identity{}
		cursor.Decode(&row)
		Identitys = append(Identitys, row)
	}

	return Identitys, nil
}
func (repo *mongoIdentityRepository) GetByID(ctx context.Context, id string) (entity.Identity, error) {
	Identity := entity.Identity{}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Identity, errors.New("Invalid Identity ID")
	}
	err = repo.Collection.FindOne(ctx, bson.M{"_id": _id}).Decode(&Identity)
	if err != nil {
		repo.logger.Errorf("FindOne transaction error: %s", err)
		return Identity, err
	}
	return Identity, nil
}
func (repo *mongoIdentityRepository) GetByDigit(ctx context.Context, digit int64) (entity.Identity, error) {
	Identity := entity.Identity{}

	err := repo.Collection.FindOne(ctx, bson.M{"digit": digit}).Decode(&Identity)

	if err != nil {
		repo.logger.Errorf("GetByEmail transaction error: %s", err)
		return Identity, err
	}

	return Identity, nil
}
func (repo *mongoIdentityRepository) GetByEmail(ctx context.Context, email string) (entity.Identity, error) {
	Identity := entity.Identity{}

	err := repo.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&Identity)

	if err != nil {
		repo.logger.Errorf("GetByEmail transaction error: %s", err)
		return Identity, err
	}

	return Identity, nil
}
func (repo *mongoIdentityRepository) Update(ctx context.Context, id string, data interface{}) (entity.Identity, error) {
	var result entity.Identity
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, errors.New("Invalid Identity ID")
	}

	Identity, err := repo.GetByID(ctx, id)
	Identity.UpdatedAt = time.Now()

	if err != nil {
		return result, err
	}
	var exist map[string]interface{}
	b, err := json.Marshal(Identity)
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
func (repo *mongoIdentityRepository) Create(ctx context.Context, Identity entity.Identity) (entity.Identity, error) {
	newIdentity := entity.Identity{}
	Identity.CreatedAt = time.Now()
	res, err := repo.Collection.InsertOne(ctx, Identity)
	id := res.InsertedID.(primitive.ObjectID).Hex()

	if err != nil {
		repo.logger.Errorf("Create transaction error: %s", err)
		return newIdentity, err
	}
	return repo.GetByID(ctx, id)
}
func (repo *mongoIdentityRepository) Delete(ctx context.Context, id string) error {

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
