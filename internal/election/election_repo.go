package election

import (
	"context"
	"fmt"
	"time"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/internal/utils"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoElectionRepository struct {
	Collection *mongo.Collection
	logger     log.Logger
}

// NewmongoElectionRepository creates a new Election repository.
func NewMongoElectionRepository(db *mongo.Database, logger log.Logger) entity.ElectionRepository {

	collection := db.Collection("elections")
	return &mongoElectionRepository{collection, logger}
}

func mongoElectionPipeline(match bson.M) []bson.M {
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
		{"$unwind": "$candidates"},
		{
			"$lookup": bson.M{
				"from":         "political_parties",
				"localField":   "candidates.political_party",
				"foreignField": "_id",
				"as":           "candidates.political_party",
			},
		},
		{"$unwind": "$candidates.political_party"},
		{"$group": bson.M{
			"_id":             "$_id",
			"description":     bson.M{"$first": "$description"},
			"title":           bson.M{"$first": "$title"},
			"pubkey":          bson.M{"$first": "$pubkey"},
			"tx_out_ref":      bson.M{"$first": "$tx_out_ref"},
			"country":         bson.M{"$first": "$country"},
			"phase":           bson.M{"$first": "$phase"},
			"accrediation_at": bson.M{"$first": "$accrediation_at"},
			"vote_at":         bson.M{"$first": "$vote_at"},
			"created_at":      bson.M{"$first": "$created_at"},
			"updated_at":      bson.M{"$first": "$updated_at"},
			"candidates":      bson.M{"$push": "$candidates"},
		},
		},
	}
}

// Fetch returns the election with the specified filter from Mongo.
func (repo *mongoElectionRepository) Fetch(ctx context.Context, filter map[string]interface{}) (res []*entity.ElectionRead, err error) {
	_filter, _ := utils.ConstructQuery(filter)

	cursor, err := repo.Collection.Aggregate(ctx, mongoElectionPipeline(_filter))

	if err != nil {
		repo.logger.Errorf("Fetch transaction error: %s", err)
		return
	}

	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("Fetch transaction error: %s", err)
		return
	}

	if len(res) == 0 {
		return []*entity.ElectionRead{}, nil
	}

	return
}

// GetByID gets the Election with the specified Id from mongo.
func (repo *mongoElectionRepository) GetByID(ctx context.Context, id string) (entity.ElectionRead, error) {
	Election := entity.ElectionRead{}

	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return Election, entity.ErrInvalidId
	}

	cursor, err := repo.Collection.Aggregate(
		ctx,
		mongoElectionPipeline(bson.M{"_id": _id}),
	)

	if err != nil {
		repo.logger.Errorf("GetByID transaction error: %s", err)
		return Election, err
	}
	var res []entity.ElectionRead

	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("GetByID transaction error: %s", err)
		return Election, err
	}
	if len(res) == 0 {
		return Election, entity.ErrNotFound
	}
	return res[0], nil
}

// GetByCountry gets the election with the specified country from mongo.
func (repo *mongoElectionRepository) GetByCountry(ctx context.Context, country string) (entity.ElectionRead, error) {
	Election := entity.ElectionRead{}

	_id, err := primitive.ObjectIDFromHex(country)

	if err != nil {
		return Election, entity.ErrInvalidId
	}

	cursor, err := repo.Collection.Aggregate(ctx, mongoElectionPipeline(bson.M{"country": _id}))

	if err != nil {
		repo.logger.Errorf("GetByCountry transaction error: %s", err)
		return Election, err
	}
	var res []entity.ElectionRead
	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("GetByCountry transaction error: %s", err)
		return Election, err
	}
	if len(res) == 0 {
		return Election, entity.ErrNotFound
	}
	return res[0], nil
}

// GetBySlug gets the Election with the specified slug from mongo.
func (repo *mongoElectionRepository) GetBySlug(ctx context.Context, slug string) (entity.ElectionRead, error) {
	Election := entity.ElectionRead{}

	cursor, err := repo.Collection.Aggregate(ctx, mongoElectionPipeline(bson.M{"slug": slug}))

	if err != nil {
		repo.logger.Errorf("GetBySlug transaction error: %s", err)
		return Election, err
	}
	var res []entity.ElectionRead
	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("GetBySlug transaction error: %s", err)
		return Election, err
	}
	if len(res) == 0 {
		return Election, entity.ErrNotFound
	}
	return res[0], nil
}

// GetWithExclude gets the Election with the specified Election excluding some other data.
func (repo *mongoElectionRepository) GetWithExclude(ctx context.Context, filter map[string]interface{}, exclude map[string]interface{}) (entity.Election, error) {
	election := entity.Election{}
	query, err := utils.ConstructQueryWithTypes(filter, map[string][]string{
		"object_id": {"country", "_id"},
	})
	query, err = utils.ConstructNotEqualQuery(exclude, query)

	if err != nil {
		switch err {
		case entity.ErrInvalidId:
			return election, entity.ErrInvalidId
		default:
			return election, err
		}
	}
	err = repo.Collection.FindOne(ctx, query).Decode(&election)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return election, entity.ErrNotFound
		default:
			repo.logger.Errorf("GetWithExclude transaction error: %s", err)
			return election, err
		}
	}

	return election, nil
}

// Get gets the Election with the specified Election
func (repo *mongoElectionRepository) Get(ctx context.Context, filter map[string]interface{}) (entity.Election, error) {
	Election := entity.Election{}
	query, err := utils.ConstructQueryWithTypes(filter, map[string][]string{
		"object_id": {"country", "_id"},
	})

	if err != nil {
		switch err {
		case entity.ErrInvalidId:
			return Election, entity.ErrInvalidId
		default:
			return Election, err
		}
	}
	err = repo.Collection.FindOne(ctx, query).Decode(&Election)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return Election, entity.ErrNotFound
		default:
			repo.logger.Errorf("GetWithExclude transaction error: %s", err)
			return Election, err
		}
	}

	return Election, nil
}

// Update updates the Election with the specified Id from mongo.
func (repo *mongoElectionRepository) Update(ctx context.Context, id string, updateElection map[string]interface{}) (entity.ElectionRead, error) {
	var result entity.ElectionRead
	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return result, entity.ErrInvalidId
	}

	election, err := repo.GetByID(ctx, id)
	if err != nil {
		return result, err
	}
	if value, ok := updateElection["country"]; ok && value != "" {
		_, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", value))
		if err != nil {
			return election, entity.ErrInvalidId
		}
		updateElection["country"] = _id
	}

	updateElection["updated_at"] = time.Now()

	_, err = repo.Collection.UpdateOne(ctx, bson.M{"_id": _id}, bson.M{"$set": updateElection})
	if err != nil {
		return result, err
	}

	return repo.GetByID(ctx, id)
}

// Create creates a new Election.
func (repo *mongoElectionRepository) Create(ctx context.Context, election entity.Election) (entity.ElectionRead, error) {
	newElection := entity.ElectionRead{}
	election.CreatedAt = time.Now()
	election.Phase = "initial"
	for k := range election.Candidates {
		election.Candidates[k].ID = primitive.NewObjectID()
	}
	res, err := repo.Collection.InsertOne(ctx, election)
	id := res.InsertedID.(primitive.ObjectID).Hex()

	if err != nil {
		repo.logger.Errorf("Create transaction error: %s", err)
		return newElection, err
	}

	return repo.GetByID(ctx, id)
}

// Delete deletes a Election.
func (repo *mongoElectionRepository) Delete(ctx context.Context, id string) error {

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
