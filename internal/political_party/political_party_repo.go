package political_party

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

type mongoPartyRepository struct {
	Collection *mongo.Collection
	logger     log.Logger
}

func NewMongoPartyRepository(db *mongo.Database, logger log.Logger) entity.PoliticalPartyRepository {

	collection := db.Collection("political_parties")
	return &mongoPartyRepository{collection, logger}
}

func (repo *mongoPartyRepository) Fetch(ctx context.Context, filter interface{}) (res []entity.PoliticalParty, err error) {
	PoliticalParties := []entity.PoliticalParty{}
	if filter == nil {
		filter = bson.M{}
	}

	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		repo.logger.Errorf("Fetch transaction error: %s", err)
		return PoliticalParties, err
	}

	for cursor.Next(ctx) {
		row := entity.PoliticalParty{}
		cursor.Decode(&row)
		PoliticalParties = append(PoliticalParties, row)
	}

	return PoliticalParties, nil
}
func (repo *mongoPartyRepository) GetByID(ctx context.Context, id string) (entity.PoliticalParty, error) {
	PoliticalParty := entity.PoliticalParty{}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return PoliticalParty, errors.New("Invalid PoliticalParty ID")
	}
	err = repo.Collection.FindOne(ctx, bson.M{"_id": _id}).Decode(&PoliticalParty)
	if err != nil {
		repo.logger.Errorf("GetByID transaction error: %s", err)
		return PoliticalParty, err
	}
	return PoliticalParty, nil
}

func (repo *mongoPartyRepository) GetByCountry(ctx context.Context, country_id string) (entity.PoliticalParty, error) {
	PoliticalParty := entity.PoliticalParty{}

	_id, err := primitive.ObjectIDFromHex(country_id)
	if err != nil {
		return PoliticalParty, errors.New("Invalid PoliticalParty ID")
	}
	err = repo.Collection.FindOne(ctx, bson.M{"country": _id}).Decode(&PoliticalParty)
	if err != nil {
		repo.logger.Errorf("GetByCountry transaction error: %s", err)
		return PoliticalParty, err
	}
	return PoliticalParty, nil
}

func (repo *mongoPartyRepository) GetBySlug(ctx context.Context, slug string) (entity.PoliticalParty, error) {
	PoliticalParty := entity.PoliticalParty{}

	err := repo.Collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&PoliticalParty)

	if err != nil {
		repo.logger.Errorf("GetBySlug transaction error: %s", err)
		return PoliticalParty, err
	}

	return PoliticalParty, nil
}
func (repo *mongoPartyRepository) Update(ctx context.Context, id string, data interface{}) (entity.PoliticalParty, error) {
	var result entity.PoliticalParty
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, errors.New("Invalid PoliticalParty ID")
	}

	PoliticalParty, err := repo.GetByID(ctx, id)

	if err != nil {
		return result, err
	}
	var exist map[string]interface{}
	b, err := json.Marshal(PoliticalParty)
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
func (repo *mongoPartyRepository) Store(ctx context.Context, PoliticalParty entity.PoliticalParty) (entity.PoliticalParty, error) {
	newPoliticalParty := entity.PoliticalParty{}
	PoliticalParty.CreatedAt = time.Now()
	res, err := repo.Collection.InsertOne(ctx, PoliticalParty)
	id := res.InsertedID.(primitive.ObjectID).Hex()

	if err != nil {
		repo.logger.Errorf("Store transaction error: %s", err)
		return newPoliticalParty, err
	}
	return repo.GetByID(ctx, id)
}
func (repo *mongoPartyRepository) Delete(ctx context.Context, id string) error {

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
