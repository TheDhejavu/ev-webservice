package politicalparty

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

type mongoPoliticalPartyRepository struct {
	Collection *mongo.Collection
	logger     log.Logger
}

// NewmongoPoliticalPartyRepository creates a new PoliticalParty repository.
func NewMongoPoliticalPartyRepository(db *mongo.Database, logger log.Logger) entity.PoliticalPartyRepository {

	collection := db.Collection("political_parties")
	return &mongoPoliticalPartyRepository{collection, logger}
}

func mongoPartyPipeline(match bson.D) mongo.Pipeline {
	matchStage := bson.D{{"$match", match}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "countries"}, {"localField", "country"}, {"foreignField", "_id"}, {"as", "country"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$country"}, {"preserveNullAndEmptyArrays", false}}}}

	return mongo.Pipeline{lookupStage, matchStage, unwindStage}
}

// Fetch returns the political parties with the specified filter from Mongo.
func (repo *mongoPoliticalPartyRepository) Fetch(ctx context.Context, filter interface{}) (res []entity.PoliticalPartyRead, err error) {
	_filter := bson.D{}

	cursor, err := repo.Collection.Aggregate(ctx, mongoPartyPipeline(_filter))

	if err != nil {
		repo.logger.Errorf("Fetch transaction error: %s", err)
		return
	}

	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("Fetch transaction error: %s", err)
		return
	}

	return
}

// GetByID gets the PoliticalParty with the specified Id from mongo.
func (repo *mongoPoliticalPartyRepository) GetByID(ctx context.Context, id string) (entity.PoliticalPartyRead, error) {
	politicalParty := entity.PoliticalPartyRead{}

	_id, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return politicalParty, entity.ErrInvalidId
	}

	cursor, err := repo.Collection.Aggregate(
		ctx,
		mongoPartyPipeline(bson.D{{"_id", _id}}),
	)

	if err != nil {
		repo.logger.Errorf("GetByID transaction error: %s", err)
		return politicalParty, err
	}
	var res []entity.PoliticalPartyRead
	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("GetByID transaction error: %s", err)
		return politicalParty, err
	}
	if len(res) == 0 {
		return politicalParty, entity.ErrNotFound
	}
	return res[0], nil
}

// GetByCountry gets the political parties with the specified country from mongo.
func (repo *mongoPoliticalPartyRepository) GetByCountry(ctx context.Context, country string) (entity.PoliticalPartyRead, error) {
	politicalParty := entity.PoliticalPartyRead{}

	_id, err := primitive.ObjectIDFromHex(country)

	if err != nil {
		return politicalParty, entity.ErrInvalidId
	}

	cursor, err := repo.Collection.Aggregate(ctx, mongoPartyPipeline(bson.D{{"country", _id}}))

	if err != nil {
		repo.logger.Errorf("GetByCountry transaction error: %s", err)
		return politicalParty, err
	}
	var res []entity.PoliticalPartyRead
	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("GetByCountry transaction error: %s", err)
		return politicalParty, err
	}
	if len(res) == 0 {
		return politicalParty, entity.ErrNotFound
	}
	return res[0], nil
}

// GetBySlug gets the PoliticalParty with the specified slug from mongo.
func (repo *mongoPoliticalPartyRepository) GetBySlug(ctx context.Context, slug string) (entity.PoliticalPartyRead, error) {
	politicalParty := entity.PoliticalPartyRead{}

	cursor, err := repo.Collection.Aggregate(ctx, mongoPartyPipeline(bson.D{{"slug", slug}}))

	if err != nil {
		repo.logger.Errorf("GetBySlug transaction error: %s", err)
		return politicalParty, err
	}
	var res []entity.PoliticalPartyRead
	if err = cursor.All(ctx, &res); err != nil {
		repo.logger.Errorf("GetBySlug transaction error: %s", err)
		return politicalParty, err
	}
	if len(res) == 0 {
		return politicalParty, entity.ErrNotFound
	}
	return res[0], nil
}

// GetWithExclude gets the PoliticalParty with the specified PoliticalParty excluding some other data.
func (repo *mongoPoliticalPartyRepository) GetWithExclude(ctx context.Context, politicalParty map[string]interface{}, exclude map[string]interface{}) (entity.PoliticalParty, error) {
	PoliticalParty := entity.PoliticalParty{}
	query := bson.M{}
	if value, ok := politicalParty["country"]; ok {
		_id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", value))
		if err != nil {
			return PoliticalParty, entity.ErrInvalidId
		}
		query, err = utils.ConstructQuery(bson.M{"country": _id})
		delete(politicalParty, "country")
	}
	query, err := utils.ConstructQuery(politicalParty)
	query, err = utils.ConstructNotEqualQuery(exclude, query)

	if err != nil {
		switch err {
		case entity.ErrInvalidId:
			return PoliticalParty, entity.ErrInvalidId
		default:
			return PoliticalParty, err
		}
	}
	err = repo.Collection.FindOne(ctx, query).Decode(&PoliticalParty)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return PoliticalParty, entity.ErrNotFound
		default:
			repo.logger.Errorf("GetWithExclude transaction error: %s", err)
			return PoliticalParty, err
		}
	}

	return PoliticalParty, nil
}

// Get gets the PoliticalParty with the specified PoliticalParty
func (repo *mongoPoliticalPartyRepository) Get(ctx context.Context, filter map[string]interface{}) (entity.PoliticalParty, error) {
	PoliticalParty := entity.PoliticalParty{}
	query := bson.M{}
	if value, ok := filter["country"]; ok {
		_id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", value))
		if err != nil {
			return PoliticalParty, entity.ErrInvalidId
		}
		query, err = utils.ConstructQuery(bson.M{"country": _id})
		delete(filter, "country")
	}
	query, err := utils.ConstructQuery(filter)

	if err != nil {
		switch err {
		case entity.ErrInvalidId:
			return PoliticalParty, entity.ErrInvalidId
		default:
			return PoliticalParty, err
		}
	}
	err = repo.Collection.FindOne(ctx, query).Decode(&PoliticalParty)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return PoliticalParty, entity.ErrNotFound
		default:
			repo.logger.Errorf("GetWithExclude transaction error: %s", err)
			return PoliticalParty, err
		}
	}

	return PoliticalParty, nil
}

// Update updates the PoliticalParty with the specified Id from mongo.
func (repo *mongoPoliticalPartyRepository) Update(ctx context.Context, id string, updatePoliticalParty map[string]interface{}) (entity.PoliticalPartyRead, error) {
	var result entity.PoliticalPartyRead
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, entity.ErrInvalidId
	}

	PoliticalParty, err := repo.GetByID(ctx, id)
	PoliticalParty.UpdatedAt = time.Now()

	if err != nil {
		return result, err
	}
	var existingPoliticalParty map[string]interface{}
	b, err := json.Marshal(PoliticalParty)
	if err != nil {
		return result, err
	}

	json.Unmarshal(b, &existingPoliticalParty)

	if value, ok := updatePoliticalParty["country"]; ok {
		_id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", value))
		if err != nil {
			return PoliticalParty, entity.ErrInvalidId
		}
		updatePoliticalParty["country"] = _id
	}

	for k := range updatePoliticalParty {
		if _, ok := existingPoliticalParty[k]; ok {
			existingPoliticalParty[k] = updatePoliticalParty[k]
		}
	}

	delete(existingPoliticalParty, "_id")
	delete(existingPoliticalParty, "createdAt")

	fmt.Println(existingPoliticalParty)

	_, err = repo.Collection.UpdateOne(ctx, bson.M{"_id": _id}, bson.M{"$set": existingPoliticalParty})
	if err != nil {
		return result, err
	}

	return repo.GetByID(ctx, id)
}

// Create creates a new PoliticalParty.
func (repo *mongoPoliticalPartyRepository) Store(ctx context.Context, PoliticalParty entity.PoliticalParty) (entity.PoliticalPartyRead, error) {
	newPoliticalParty := entity.PoliticalPartyRead{}
	PoliticalParty.CreatedAt = time.Now()
	res, err := repo.Collection.InsertOne(ctx, PoliticalParty)
	id := res.InsertedID.(primitive.ObjectID).Hex()

	if err != nil {
		repo.logger.Errorf("Create transaction error: %s", err)
		return newPoliticalParty, err
	}
	return repo.GetByID(ctx, id)
}

// Delete deletes a PoliticalParty.
func (repo *mongoPoliticalPartyRepository) Delete(ctx context.Context, id string) error {

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
