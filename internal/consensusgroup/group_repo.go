package consensusgroup

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/internal/utils"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoGroupRepository struct {
	Collection *mongo.Collection
	logger     log.Logger
}

// NewmongoGroupRepository creates a new ConsensusGroup repository.
func NewMongoGroupRepository(db *mongo.Database, logger log.Logger) entity.ConsensusGroupRepository {

	collection := db.Collection("consensus_group")
	return &mongoGroupRepository{collection, logger}
}

func mongoPartyPipeline(match bson.M) []bson.M {
	return []bson.M{
		{
			"$match": match,
		},
		{
			"$lookup": bson.M{
				"from":         "countries",
				"localField":   "country",
				"foreignField": "_id",
				"as":           "country",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$country",
				"preserveNullAndEmptyArrays": false,
			},
		},
	}
}

// Fetch returns the political parties with the specified filter from Mongo.
func (repo *mongoGroupRepository) Fetch(ctx context.Context, filter interface{}) (res []entity.ConsensusGroupRead, err error) {
	_filter := bson.M{}

	cursor, err := repo.Collection.Aggregate(ctx, mongoPartyPipeline(_filter))

	if err != nil {
		repo.logger.Errorf("Fetch transaction error: %s", err)
		return
	}

	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("Fetch transaction error: %s", err)
		return
	}

	if len(res) == 0 {
		res = []entity.ConsensusGroupRead{}
		return
	}

	return
}

// GetByID gets the ConsensusGroup with the specified Id from mongo.
func (repo *mongoGroupRepository) GetByID(ctx context.Context, id string) (entity.ConsensusGroupRead, error) {
	ConsensusGroup := entity.ConsensusGroupRead{}

	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return ConsensusGroup, entity.ErrInvalidId
	}

	cursor, err := repo.Collection.Aggregate(
		ctx,
		mongoPartyPipeline(bson.M{"_id": _id}),
	)

	if err != nil {
		repo.logger.Errorf("GetByID transaction error: %s", err)
		return ConsensusGroup, err
	}
	var res []entity.ConsensusGroupRead
	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("GetByID transaction error: %s", err)
		return ConsensusGroup, err
	}
	if len(res) == 0 {
		return ConsensusGroup, entity.ErrNotFound
	}
	return res[0], nil
}

// GetByCountry gets the political parties with the specified country from mongo.
func (repo *mongoGroupRepository) GetByCountry(ctx context.Context, country string) (entity.ConsensusGroupRead, error) {
	ConsensusGroup := entity.ConsensusGroupRead{}

	_id, err := primitive.ObjectIDFromHex(country)

	if err != nil {
		return ConsensusGroup, entity.ErrInvalidId
	}

	cursor, err := repo.Collection.Aggregate(ctx, mongoPartyPipeline(bson.M{"country": _id}))

	if err != nil {
		repo.logger.Errorf("GetByCountry transaction error: %s", err)
		return ConsensusGroup, err
	}
	var res []entity.ConsensusGroupRead
	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("GetByCountry transaction error: %s", err)
		return ConsensusGroup, err
	}
	if len(res) == 0 {
		return ConsensusGroup, entity.ErrNotFound
	}
	return res[0], nil
}

// GetByPubkey gets the ConsensusGroup with the specified slug from mongo.
func (repo *mongoGroupRepository) GetByPubKey(ctx context.Context, publicKey string) (entity.ConsensusGroupRead, error) {
	ConsensusGroup := entity.ConsensusGroupRead{}

	cursor, err := repo.Collection.Aggregate(ctx, mongoPartyPipeline(bson.M{"public_key": publicKey}))

	if err != nil {
		repo.logger.Errorf("GetByPubkey transaction error: %s", err)
		return ConsensusGroup, err
	}
	var res []entity.ConsensusGroupRead
	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("GetByPubkey transaction error: %s", err)
		return ConsensusGroup, err
	}
	if len(res) == 0 {
		return ConsensusGroup, entity.ErrNotFound
	}
	return res[0], nil
}

// GetWithExclude gets the ConsensusGroup with the specified ConsensusGroup excluding some other data.
func (repo *mongoGroupRepository) GetWithExclude(ctx context.Context, group map[string]interface{}, exclude map[string]interface{}) (entity.ConsensusGroup, error) {
	ConsensusGroup := entity.ConsensusGroup{}
	query := bson.M{}
	if value, ok := group["country"]; ok {

		_id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", value))
		if err != nil {
			return ConsensusGroup, entity.ErrInvalidId
		}
		query, err = utils.ConstructQuery(bson.M{"country": _id})
		delete(group, "country")
	}
	query, err := utils.ConstructQuery(group)
	query, err = utils.ConstructNotEqualQuery(exclude, query)

	if err != nil {
		switch err {
		case entity.ErrInvalidId:
			return ConsensusGroup, entity.ErrInvalidId
		default:
			return ConsensusGroup, err
		}
	}
	err = repo.Collection.FindOne(ctx, query).Decode(&ConsensusGroup)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return ConsensusGroup, entity.ErrNotFound
		default:
			repo.logger.Errorf("GetWithExclude transaction error: %s", err)
			return ConsensusGroup, err
		}
	}
	return ConsensusGroup, nil
}

// Get gets the ConsensusGroup with the specified ConsensusGroup
func (repo *mongoGroupRepository) Get(ctx context.Context, filter map[string]interface{}) (entity.ConsensusGroupRead, error) {

	ConsensusGroup := entity.ConsensusGroupRead{}
	query, err := utils.ConstructQueryWithTypes(filter, map[string][]string{
		"object_id": {"country", "_id"},
	})

	if err != nil {
		return ConsensusGroup, err
	}
	err = repo.Collection.FindOne(ctx, query).Decode(&ConsensusGroup)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return ConsensusGroup, entity.ErrNotFound
		default:
			repo.logger.Errorf("GetWithExclude transaction error: %s", err)
			return ConsensusGroup, err
		}
	}

	return ConsensusGroup, nil
}

// Update updates the ConsensusGroup with the specified Id from mongo.
func (repo *mongoGroupRepository) Update(ctx context.Context, id string, updateConsensusGroup map[string]interface{}) (entity.ConsensusGroupRead, error) {
	var result entity.ConsensusGroupRead
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, entity.ErrInvalidId
	}

	ConsensusGroup, err := repo.GetByID(ctx, id)
	ConsensusGroup.UpdatedAt = time.Now()

	if err != nil {
		return result, err
	}
	var existingConsensusGroup map[string]interface{}
	b, err := json.Marshal(ConsensusGroup)
	if err != nil {
		return result, err
	}

	json.Unmarshal(b, &existingConsensusGroup)

	if value, ok := updateConsensusGroup["country"]; ok {
		_id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", value))
		if err != nil {
			return ConsensusGroup, entity.ErrInvalidId
		}
		updateConsensusGroup["country"] = _id
	}

	for k := range updateConsensusGroup {
		if _, ok := existingConsensusGroup[k]; ok {
			existingConsensusGroup[k] = updateConsensusGroup[k]
		}
	}

	delete(existingConsensusGroup, "_id")
	delete(existingConsensusGroup, "createdAt")

	fmt.Println(existingConsensusGroup)

	_, err = repo.Collection.UpdateOne(ctx, bson.M{"_id": _id}, bson.M{"$set": existingConsensusGroup})
	if err != nil {
		return result, err
	}

	return repo.GetByID(ctx, id)
}

// Create creates a new ConsensusGroup.
func (repo *mongoGroupRepository) Create(ctx context.Context, ConsensusGroup entity.ConsensusGroup) (entity.ConsensusGroupRead, error) {
	newConsensusGroup := entity.ConsensusGroupRead{}
	ConsensusGroup.CreatedAt = time.Now()
	res, err := repo.Collection.InsertOne(ctx, ConsensusGroup)
	id := res.InsertedID.(primitive.ObjectID).Hex()

	if err != nil {
		repo.logger.Errorf("Create transaction error: %s", err)
		return newConsensusGroup, err
	}
	return repo.GetByID(ctx, id)
}

// Delete deletes a ConsensusGroup.
func (repo *mongoGroupRepository) Delete(ctx context.Context, id string) error {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repo.logger.Errorf("Delete transaction error: %s", err)
		return err
	}
	_, err = repo.Collection.DeleteOne(ctx, bson.M{"_id": _id})

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return entity.ErrNotFound
		default:
			repo.logger.Errorf("Delete transaction error: %s", err)
			return err
		}
	}

	return err
}
