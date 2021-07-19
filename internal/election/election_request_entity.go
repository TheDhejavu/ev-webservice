package election

import (
	"context"
	"fmt"
	"time"

	validator "github.com/go-playground/validator/v10"
)

type createElectionRequest struct {
	Title       string `json:"title"  validate:"required"`
	Description string `json:"description"  validate:"required"`
	Country     string `json:"country"  validate:"required,not_exists"`
	Candidates  []*struct {
		FullName       string `json:"full_name" bson:"title,omitempty" validate:"required"`
		Position       string `json:"position" bson:"position,omitempty"  validate:"required"`
		PoliticalParty string `json:"political_party" bson:"political_party,omitempty" validate:"required,not_exists"`
	} `json:"candidates"  validate:"required,dive,min=1"`
	AccreditationAt createRequestAt `json:"accrediation_at" validate:"required"`
	VoteAt          createRequestAt `json:"vote_at"  validate:"required"`
}

type createRequestAt struct {
	Start time.Time `json:"start" validate:"required"`
	End   time.Time `json:"end" validate:"required"`
}

func (request createElectionRequest) Validate(ctx context.Context, handler electionHandler) error {

	handler.v.Validator.RegisterValidation("not_exists", func(fl validator.FieldLevel) bool {
		fieldName := fl.FieldName()
		if fieldName == "Country" {
			value, _ := handler.countryService.IdExists(ctx, fl.Field().String())

			if !value {
				return false
			}
			return true
		}
		fmt.Println(fieldName)
		if fieldName == "PoliticalParty" {

			value, _ := handler.partyService.Exists(ctx, map[string]interface{}{
				"_id": fl.Field().String(),
			}, nil)
			fmt.Println(value)
			if !value {
				return false
			}
			return true
		}
		return true
	})

	err := handler.v.Validator.Struct(request)

	return err
}

type updateElectionRequest map[string]interface{}

var (
	updateElectionRule = map[string]interface{}{
		// "country": "not_exists,omitempty",
		// "phase":   "oneof=accreditation voting initial,omitempty",
		// "candidates": map[string]interface{}{
		// 	"political_party": "not_exists,omitempty",
		// },
	}
)

func (request updateElectionRequest) Validate(ctx context.Context, handler electionHandler, params requestParams) error {
	handler.v.Validator.RegisterValidation("not_exists", func(fl validator.FieldLevel) bool {
		fieldName := fl.FieldName()
		fieldValue := fl.Field().String()
		if fieldName == "Country" && fieldValue != "" {
			value, _ := handler.countryService.IdExists(ctx, fieldValue)
			if !value {
				return false
			}
			return true
		}
		if fieldName == "PoliticalParty" && fieldValue != "" {
			value, _ := handler.partyService.Exists(ctx, map[string]interface{}{
				"political_party": request["country"],
			}, nil)
			if !value {
				return false
			}
			return true
		}
		return true
	})
	// fmt.Println(request)
	errs := handler.v.Validator.ValidateMap(request, updateElectionRule)
	if len(errs) > 0 {
		var errors validator.ValidationErrors
		var fieldError validator.FieldError

		// fmt.Println("VALIDATION_ERROR:", errs)
		errors = append(errors, fieldError)
		return errors
	}

	return nil
}

type requestParams struct {
	Id string `uri:"id" validate:"required"`
}
