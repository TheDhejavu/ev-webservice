package utils

import (
	"fmt"

	"github.com/workspace/evoting/ev-webservice/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConstructNotEqualQuery(data map[string]interface{}, query bson.M) (bson.M, error) {
	for field, value := range data {
		if field == "id" {
			id := fmt.Sprintf("%v", value)
			_id, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return query, entity.ErrInvalidId
			}
			query["_id"] = bson.M{"$ne": _id}
		} else {
			query[field] = bson.M{"$ne": value}
		}
	}

	return query, nil
}
func ConstructQuery(data map[string]interface{}) (bson.M, error) {
	query := bson.M{}

	for field, value := range data {

		if field == "id" {
			id := fmt.Sprintf("%v", value)
			_id, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				return query, entity.ErrInvalidId
			}
			query["_id"] = _id
		} else {
			query[field] = value
		}
	}
	return query, nil
}

func ConstructQueryWithTypes(filter map[string]interface{}, types map[string][]string) (query bson.M, err error) {
	for k, values := range types {
		if k == "object_id" {
			for i := 0; i < len(values); i++ {
				value := filter[values[i]]
				key := values[i]
				if _, ok := filter[key]; ok {
					_id, err := primitive.ObjectIDFromHex(fmt.Sprintf("%s", value))
					if err != nil {
						return query, entity.ErrInvalidId
					}
					filter[key] = _id
				}
			}
		}
	}

	query, err = ConstructQuery(filter)
	fmt.Println(query)
	return query, nil
}
