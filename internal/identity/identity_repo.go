package identity

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/workspace/evoting/ev-webservice/internal/config"
	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/internal/utils"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoIdentityRepository struct {
	Collection *mongo.Collection
	logger     log.Logger
	config     config.Config
}

// NewmongoIdentityRepository creates a new Identity repository.
func NewMongoIdentityRepository(db *mongo.Database, logger log.Logger, config config.Config) entity.IdentityRepository {

	collection := db.Collection("identities")
	return &mongoIdentityRepository{collection, logger, config}
}

func mongoIdentityPipeline(match bson.M) []bson.M {
	return []bson.M{
		{
			"$match": match,
		},
		{
			"$lookup": bson.M{
				"from":         "countries",
				"localField":   "origin.country",
				"foreignField": "_id",
				"as":           "origin.country",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$origin.country",
				"preserveNullAndEmptyArrays": false,
			},
		},
		{
			"$lookup": bson.M{
				"from":         "countries",
				"localField":   "residence.country",
				"foreignField": "_id",
				"as":           "residence.country",
			},
		},
		{
			"$unwind": bson.M{
				"path":                       "$residence.country",
				"preserveNullAndEmptyArrays": false,
			},
		},
	}
}

// Fetch returns the Identity with the specified filter from Mongo.
func (repo *mongoIdentityRepository) Fetch(ctx context.Context, filter interface{}) (res []entity.IdentityRead, err error) {
	_filter := bson.M{}

	cursor, err := repo.Collection.Aggregate(ctx, mongoIdentityPipeline(_filter))

	if err != nil {
		repo.logger.Errorf("Fetch transaction error: %s", err)
		return
	}

	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("Fetch transaction error: %s", err)
		return
	}

	if len(res) == 0 {
		return []entity.IdentityRead{}, nil
	}

	return
}

// GetByID gets the Identity with the specified Id from mongo.
func (repo *mongoIdentityRepository) GetByID(ctx context.Context, id string) (entity.IdentityRead, error) {
	Identity := entity.IdentityRead{}

	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return Identity, entity.ErrInvalidId
	}

	cursor, err := repo.Collection.Aggregate(
		ctx,
		mongoIdentityPipeline(bson.M{"_id": _id}),
	)

	if err != nil {
		repo.logger.Errorf("GetByID transaction error: %s", err)
		return Identity, err
	}
	var res []entity.IdentityRead

	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("GetByID transaction error: %s", err)
		return Identity, err
	}
	if len(res) == 0 {
		return Identity, entity.ErrNotFound
	}
	return res[0], nil
}

// GetByCountry gets the Identity with the specified country from mongo.
func (repo *mongoIdentityRepository) GetByCountry(ctx context.Context, country string) (entity.IdentityRead, error) {
	Identity := entity.IdentityRead{}

	_id, err := primitive.ObjectIDFromHex(country)

	if err != nil {
		return Identity, entity.ErrInvalidId
	}

	cursor, err := repo.Collection.Aggregate(ctx, mongoIdentityPipeline(bson.M{"country": _id}))

	if err != nil {
		repo.logger.Errorf("GetByCountry transaction error: %s", err)
		return Identity, err
	}
	var res []entity.IdentityRead
	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("GetByCountry transaction error: %s", err)
		return Identity, err
	}
	if len(res) == 0 {
		return Identity, entity.ErrNotFound
	}
	return res[0], nil
}

// GetByEmail gets the Identity with the specified country from mongo.
func (repo *mongoIdentityRepository) GetByEmail(ctx context.Context, email string) (entity.IdentityRead, error) {
	Identity := entity.IdentityRead{}

	cursor, err := repo.Collection.Aggregate(ctx, mongoIdentityPipeline(bson.M{"email": email}))

	if err != nil {
		repo.logger.Errorf("GetByCountry transaction error: %s", err)
		return Identity, err
	}
	var res []entity.IdentityRead
	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("GetByCountry transaction error: %s", err)
		return Identity, err
	}
	if len(res) == 0 {
		return Identity, entity.ErrNotFound
	}
	return res[0], nil
}

// GetByDigits gets the Identity with the specified digits from mongo.
func (repo *mongoIdentityRepository) GetByDigits(ctx context.Context, digits uint64) (entity.IdentityRead, error) {
	Identity := entity.IdentityRead{}

	cursor, err := repo.Collection.Aggregate(ctx, mongoIdentityPipeline(bson.M{"digits": digits}))

	if err != nil {
		repo.logger.Errorf("GetByDigits transaction error: %s", err)
		return Identity, err
	}
	var res []entity.IdentityRead
	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("GetByDigits transaction error: %s", err)
		return Identity, err
	}

	if len(res) == 0 {
		return Identity, entity.ErrNotFound
	}

	identity := res[0]
	identity.BirthCertificate = fmt.Sprintf("%s/%s", repo.config.AssetsURL, identity.BirthCertificate)
	identity.NationalIdCard = fmt.Sprintf("%s/%s", repo.config.AssetsURL, identity.NationalIdCard)
	identity.VoterCard = fmt.Sprintf("%s/%s", repo.config.AssetsURL, identity.VoterCard)

	return identity, nil
}

// GetWithExclude gets the Identity with the specified Identity excluding some other data.
func (repo *mongoIdentityRepository) GetWithExclude(ctx context.Context, filter map[string]interface{}, exclude map[string]interface{}) (entity.Identity, error) {
	Identity := entity.Identity{}
	query, err := utils.ConstructQueryWithTypes(filter, map[string][]string{
		"object_id": {"country", "_id"},
	})
	query, err = utils.ConstructNotEqualQuery(exclude, query)

	if err != nil {
		switch err {
		case entity.ErrInvalidId:
			return Identity, entity.ErrInvalidId
		default:
			return Identity, err
		}
	}
	err = repo.Collection.FindOne(ctx, query).Decode(&Identity)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return Identity, entity.ErrNotFound
		default:
			repo.logger.Errorf("GetWithExclude transaction error: %s", err)
			return Identity, err
		}
	}

	return Identity, nil
}

// Get gets the Identity with the specified Identity
func (repo *mongoIdentityRepository) Get(ctx context.Context, filter map[string]interface{}) (entity.Identity, error) {
	Identity := entity.Identity{}
	query, err := utils.ConstructQueryWithTypes(filter, map[string][]string{
		"object_id": {"country", "_id"},
	})

	if err != nil {
		switch err {
		case entity.ErrInvalidId:
			return Identity, entity.ErrInvalidId
		default:
			return Identity, err
		}
	}
	err = repo.Collection.FindOne(ctx, query).Decode(&Identity)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return Identity, entity.ErrNotFound
		default:
			repo.logger.Errorf("GetWithExclude transaction error: %s", err)
			return Identity, err
		}
	}

	return Identity, nil
}

// Update updates the Identity with the specified Id from mongo.
func (repo *mongoIdentityRepository) Update(ctx context.Context, id string, updateIdentity map[string]interface{}) (entity.IdentityRead, error) {
	var result entity.IdentityRead
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, entity.ErrInvalidId
	}

	Identity, err := repo.GetByID(ctx, id)
	Identity.UpdatedAt = time.Now()

	if err != nil {
		return result, err
	}
	var existingIdentity entity.Identity

	b, err := json.Marshal(Identity)
	if err != nil {
		return result, err
	}

	json.Unmarshal(b, &existingIdentity)

	var upIdentity entity.Identity
	b, err = json.Marshal(updateIdentity)
	if err != nil {
		return result, err
	}
	json.Unmarshal(b, &upIdentity)

	if value, ok := updateIdentity["country"]; ok {
		_id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", value))
		if err != nil {
			return Identity, entity.ErrInvalidId
		}
		updateIdentity["country"] = _id
	}

	_, err = repo.Collection.UpdateOne(ctx, bson.M{"_id": _id}, bson.M{"$set": existingIdentity})
	if err != nil {
		return result, err
	}

	return repo.GetByID(ctx, id)
}

// Create creates a new Identity.
func (repo *mongoIdentityRepository) Create(ctx context.Context, Identity entity.Identity) (entity.IdentityRead, error) {
	newIdentity := entity.IdentityRead{}

	Identity.CreatedAt = time.Now()
	Identity.Digits = utils.UniqueDigits()

	res, err := repo.Collection.InsertOne(ctx, Identity)
	if err != nil {
		repo.logger.Errorf("Create transaction error: %s", err)
		return newIdentity, err
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()

	if err != nil {
		repo.logger.Errorf("Create transaction error: %s", err)
		return newIdentity, err
	}

	return repo.GetByID(ctx, id)
}

// Delete deletes a Identity.
func (repo *mongoIdentityRepository) Delete(ctx context.Context, id string) error {

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
