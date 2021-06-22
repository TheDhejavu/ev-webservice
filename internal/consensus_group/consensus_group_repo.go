package consensus_group

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

type mongoGroupRepository struct {
	Collection *mongo.Collection
	logger     log.Logger
}

func NewMongoGroupRepository(db *mongo.Database, logger log.Logger) entity.ConsensusGroupRepository {

	collection := db.Collection("consensus_groups")
	return &mongoGroupRepository{collection, logger}
}

func (repo *mongoGroupRepository) Fetch(ctx context.Context, filter interface{}) (res []entity.ConsensusGroup, err error) {
	ConsensusGroups := []entity.ConsensusGroup{}
	if filter == nil {
		filter = bson.M{}
	}

	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		repo.logger.Errorf("Fetch transaction error: %s", err)
		return ConsensusGroups, err
	}

	for cursor.Next(ctx) {
		row := entity.ConsensusGroup{}
		cursor.Decode(&row)
		ConsensusGroups = append(ConsensusGroups, row)
	}

	return ConsensusGroups, nil
}
func (repo *mongoGroupRepository) GetByID(ctx context.Context, id string) (entity.ConsensusGroup, error) {
	ConsensusGroup := entity.ConsensusGroup{}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ConsensusGroup, errors.New("Invalid ConsensusGroup ID")
	}
	err = repo.Collection.FindOne(ctx, bson.M{"_id": _id}).Decode(&ConsensusGroup)
	if err != nil {
		repo.logger.Errorf("GetByID transaction error: %s", err)
		return ConsensusGroup, err
	}
	return ConsensusGroup, nil
}

func (repo *mongoGroupRepository) GetByPubKey(ctx context.Context, public_key []byte) (entity.ConsensusGroup, error) {
	ConsensusGroup := entity.ConsensusGroup{}

	err := repo.Collection.FindOne(ctx, bson.M{"country": public_key}).Decode(&ConsensusGroup)
	if err != nil {
		repo.logger.Errorf("GetByCountry transaction error: %s", err)
		return ConsensusGroup, err
	}
	return ConsensusGroup, nil
}

func (repo *mongoGroupRepository) Update(ctx context.Context, id string, data interface{}) (entity.ConsensusGroup, error) {
	var result entity.ConsensusGroup
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, errors.New("Invalid ConsensusGroup ID")
	}

	ConsensusGroup, err := repo.GetByID(ctx, id)

	if err != nil {
		return result, err
	}
	var exist map[string]interface{}
	b, err := json.Marshal(ConsensusGroup)
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
func (repo *mongoGroupRepository) Create(ctx context.Context, ConsensusGroup entity.ConsensusGroup) (entity.ConsensusGroup, error) {
	newConsensusGroup := entity.ConsensusGroup{}
	ConsensusGroup.CreatedAt = time.Now()
	res, err := repo.Collection.InsertOne(ctx, ConsensusGroup)
	id := res.InsertedID.(primitive.ObjectID).Hex()

	if err != nil {
		repo.logger.Errorf("Create transaction error: %s", err)
		return newConsensusGroup, err
	}
	return repo.GetByID(ctx, id)
}
func (repo *mongoGroupRepository) Delete(ctx context.Context, id string) error {

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
