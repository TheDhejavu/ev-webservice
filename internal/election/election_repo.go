package election

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoElectionRepository struct {
	Collection *mongo.Collection
	logger     log.Logger
}

func NewMongoElectionRepository(db *mongo.Database, logger log.Logger) entity.ElectionRepository {

	collection := db.Collection("countries")
	return &mongoElectionRepository{collection, logger}
}

func (repo *mongoElectionRepository) Fetch(ctx context.Context, filter interface{}) (res []entity.Election, err error) {
	Countries := []entity.Election{}
	if filter == nil {
		filter = bson.M{}
	}

	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		repo.logger.Errorf("Fetch transaction error: %s", err)
		return Countries, err
	}

	for cursor.Next(ctx) {
		row := entity.Election{}
		cursor.Decode(&row)
		Countries = append(Countries, row)
	}

	return Countries, nil
}
func (repo *mongoElectionRepository) GetByID(ctx context.Context, id string) (entity.Election, error) {
	Election := entity.Election{}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Election, errors.New("Invalid Election ID")
	}
	err = repo.Collection.FindOne(ctx, bson.M{"_id": _id}).Decode(&Election)
	if err != nil {
		repo.logger.Errorf("FindOne transaction error: %s", err)
		return Election, err
	}
	return Election, nil
}

func (repo *mongoElectionRepository) Update(ctx context.Context, id string, data interface{}) (entity.Election, error) {
	var result entity.Election
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, errors.New("Invalid Election ID")
	}

	Election, err := repo.GetByID(ctx, id)

	if err != nil {
		return result, err
	}
	var exist map[string]interface{}
	b, err := json.Marshal(Election)
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
func (repo *mongoElectionRepository) Create(ctx context.Context, Election entity.Election) (entity.Election, error) {
	newElection := entity.Election{}

	res, err := repo.Collection.InsertOne(ctx, Election)
	id := res.InsertedID.(primitive.ObjectID).Hex()

	if err != nil {
		repo.logger.Errorf("Create transaction error: %s", err)
		return newElection, err
	}
	return repo.GetByID(ctx, id)
}
func (repo *mongoElectionRepository) Delete(ctx context.Context, id string) error {

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
