package country

import (
	"context"
	"encoding/json"
	"time"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"github.com/workspace/evoting/ev-webservice/internal/utils"
	"github.com/workspace/evoting/ev-webservice/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoCountryRepository struct {
	Collection *mongo.Collection
	logger     log.Logger
}

// NewmongoCountryRepository creates a new Country repository.
func NewMongoCountryRepository(db *mongo.Database, logger log.Logger) entity.CountryRepository {

	collection := db.Collection("countries")
	return &mongoCountryRepository{collection, logger}
}

// Fetch returns the Countrys with the specified filter from Mongo.
func (repo *mongoCountryRepository) Fetch(ctx context.Context, filter interface{}) (res []entity.Country, err error) {
	Countries := []entity.Country{}
	if filter == nil {
		filter = bson.M{}
	}

	cursor, err := repo.Collection.Find(
		ctx,
		filter,
	)
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

// GetByID gets the Country with the specified Id from mongo.
func (repo *mongoCountryRepository) GetByID(ctx context.Context, id string) (entity.Country, error) {
	Country := entity.Country{}

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Country, entity.ErrInvalidId
	}
	err = repo.Collection.FindOne(ctx, bson.M{"_id": _id}).Decode(&Country)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return Country, entity.ErrNotFound
		default:
			repo.logger.Errorf("getById transaction error: %s", err)
			return Country, err
		}
	}
	return Country, nil
}

// GetBySlug gets the Country with the specified slug from mongo.
func (repo *mongoCountryRepository) GetBySlug(ctx context.Context, slug string) (entity.Country, error) {
	Country := entity.Country{}

	err := repo.Collection.FindOne(ctx, bson.M{"slug": slug}).Decode(&Country)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return Country, entity.ErrNotFound
		default:
			repo.logger.Errorf("GetBySlug transaction error: %s", err)
			return Country, err
		}
	}

	return Country, nil
}

// GetByName gets the Country with the specified name from mongo.
func (repo *mongoCountryRepository) GetByName(ctx context.Context, name string) (entity.Country, error) {
	Country := entity.Country{}

	err := repo.Collection.FindOne(ctx, bson.M{"name": name}).Decode(&Country)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return Country, entity.ErrNotFound
		default:
			repo.logger.Errorf("GetByName transaction error: %s", err)
			return Country, err
		}
	}

	return Country, nil
}

// GetWithExclude gets the Country with the specified Country excluding some other data.
func (repo *mongoCountryRepository) GetWithExclude(ctx context.Context, country map[string]interface{}, exclude map[string]interface{}) (entity.Country, error) {
	Country := entity.Country{}
	query, err := utils.ConstructQuery(country)
	query, err = utils.ConstructNotEqualQuery(exclude, query)

	if err != nil {
		switch err {
		case entity.ErrInvalidId:
			return Country, entity.ErrInvalidId
		default:
			return Country, err
		}
	}
	err = repo.Collection.FindOne(ctx, query).Decode(&Country)

	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return Country, entity.ErrNotFound
		default:
			repo.logger.Errorf("GetWithExclude transaction error: %s", err)
			return Country, err
		}
	}

	return Country, nil
}

// Update updates the Country with the specified Id from mongo.
func (repo *mongoCountryRepository) Update(ctx context.Context, id string, updateCountry map[string]interface{}) (entity.Country, error) {
	var result entity.Country
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return result, entity.ErrInvalidId
	}

	Country, err := repo.GetByID(ctx, id)
	Country.UpdatedAt = time.Now()

	if err != nil {
		return result, err
	}
	var existingCountry map[string]interface{}
	b, err := json.Marshal(Country)
	if err != nil {
		return result, err
	}

	json.Unmarshal(b, &existingCountry)

	for k := range updateCountry {
		if _, ok := existingCountry[k]; ok {
			existingCountry[k] = updateCountry[k]
		}
	}

	delete(existingCountry, "_id")
	delete(existingCountry, "id")
	delete(existingCountry, "createdAt")

	_, err = repo.Collection.UpdateOne(ctx, bson.M{"_id": _id}, bson.M{"$set": existingCountry})
	if err != nil {
		return result, err
	}

	return repo.GetByID(ctx, id)
}

// Create creates a new Country.
func (repo *mongoCountryRepository) Store(ctx context.Context, Country entity.Country) (entity.Country, error) {
	newCountry := entity.Country{}
	Country.CreatedAt = time.Now()
	res, err := repo.Collection.InsertOne(ctx, Country)
	id := res.InsertedID.(primitive.ObjectID).Hex()

	if err != nil {
		repo.logger.Errorf("Create transaction error: %s", err)
		return newCountry, err
	}
	return repo.GetByID(ctx, id)
}

// Delete deletes a Country.
func (repo *mongoCountryRepository) Delete(ctx context.Context, id string) error {

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		repo.logger.Errorf("DeleteOne transaction error: %s", err)
		return entity.ErrInvalidId
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
