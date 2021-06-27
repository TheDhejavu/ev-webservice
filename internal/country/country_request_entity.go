package country

import (
	"context"

	validator "github.com/go-playground/validator/v10"
)

type createCountryRequest struct {
	Flag string `json:"flag" validate:"required"`
	Name string `json:"name" validate:"required,exists"`
	Slug string `json:"slug" validate:"required,exists"`
}

func (request createCountryRequest) Validate(ctx context.Context, handler countryHandler) error {
	handler.v.Validator.RegisterValidation("exists", func(fl validator.FieldLevel) bool {
		fieldName := fl.FieldName()
		if fieldName == "Slug" {
			value, _ := handler.service.SlugExists(ctx, fl.Field().String(), nil)
			if value {
				return false
			}
			return true
		}
		if fieldName == "Name" {
			value, _ := handler.service.NameExists(ctx, fl.Field().String(), nil)
			if value {
				return false
			}
			return true
		}
		return true
	})

	err := handler.v.Validator.Struct(request)

	return err
}

type updateCountryRequest struct {
	Flag string `json:"flag"`
	Name string `json:"name" validate:"exists"`
	Slug string `json:"slug" validate:"exists"`
}

func (request updateCountryRequest) Validate(ctx context.Context, handler countryHandler, params countryRequestParams) error {
	handler.v.Validator.RegisterValidation("exists", func(fl validator.FieldLevel) bool {
		fieldName := fl.FieldName()
		fieldValue := fl.Field().String()
		if fieldName == "Slug" && fieldValue != "" {
			value, _ := handler.service.SlugExists(ctx, fl.Field().String(), params.Id)
			if value {
				return false
			}
			return true
		}
		if fieldName == "Name" && fieldValue != "" {
			value, _ := handler.service.NameExists(ctx, fl.Field().String(), params.Id)
			if value {
				return false
			}
			return true
		}
		return true
	})

	err := handler.v.Validator.Struct(request)

	return err
}

type countryRequestParams struct {
	Id string `uri:"id" validate:"required"`
}
