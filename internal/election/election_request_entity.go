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

type updateElectionRequest struct {
	Title       string `json:"title" `
	Description string `json:"description"`
	Country     string `json:"country" validate:"not_exists"`
	Phase       string `json:"phase" bson:"phase,omitempty" validate:"oneof=accreditation voting start"`
	Candidates  []struct {
		FullName       string `json:"full_name" `
		Position       string `json:"position"`
		PoliticalParty string `json:"political_party"  validate:"not_exists"`
	} `json:"candidates" validate:"dive"`
	AccreditationAt updateRequestAt `json:"accreditation_at"`
	VoteAt          updateRequestAt `json:"vote_at"`
}

type updateRequestAt struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

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
				"political_party": request.Country,
			}, nil)
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

type requestParams struct {
	Id string `uri:"id" validate:"required"`
}
