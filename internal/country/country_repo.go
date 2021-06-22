package country

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

type mongoCountryRepository struct {
	Collection *mongo.Collection
	logger     log.Logger
}

func NewMongoCountryRepository(db *mongo.Database, logger log.Logger) entity.CountryRepository {

	collection := db.Collection("countries")
	return &mongoCountryRepository{collection, logger}
}

func (repo *mongoCountryRepository) Fetch(ctx context.Context, filter interface{}) (res []entity.Country, err error) {
	Countries := []entity.Country{}
	if filter == nil {
		filter = bson.M{}
	}

	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		repo.logger.Errorf("Fetch transaction error: %s", err)
		return Countries, err
	}

	for cursor.Next(ctx) {
		row := entity.Country{}
		cursor.Decode(&row)
		Countries = append(Countries, row)
	}

	return Countries, nil
}
func (repo *mongoCountryRepository) GetByID(ctx context.Context, id string) (entity.Country, error) {
	Country := entity.Country{}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Country, errors.New("Invalid Country ID")
	}
	err = repo.Collection.FindOne(ctx, bson.M{"_id": _id}).Decode(&Country)
	if err != nil {
		repo.logger.Errorf("FindOne transaction error: %s", err)
		return Country, err
	}
	return Country, nil
}

func (repo *mongoCountryRepository) Update(ctx context.Context, id string, data interface{}) (entity.Country, error) {
	var result entity.Country
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, errors.New("Invalid Country ID")
	}

	Country, err := repo.GetByID(ctx, id)

	if err != nil {
		return result, err
	}
	var exist map[string]interface{}
	b, err := json.Marshal(Country)
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
func (repo *mongoCountryRepository) Store(ctx context.Context, Country entity.Country) (entity.Country, error) {
	newCountry := entity.Country{}

	res, err := repo.Collection.InsertOne(ctx, Country)
	id := res.InsertedID.(primitive.ObjectID).Hex()

	if err != nil {
		repo.logger.Errorf("Store transaction error: %s", err)
		return newCountry, err
	}
	return repo.GetByID(ctx, id)
}
func (repo *mongoCountryRepository) Delete(ctx context.Context, id string) error {

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
