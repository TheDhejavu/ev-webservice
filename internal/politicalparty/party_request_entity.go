package politicalparty

import (
	"context"

	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
)

type createPoliticalPartyRequest struct {
	Name    string `json:"name" validate:"required,exists"`
	Slug    string `json:"slug" validate:"required,exists"`
	Country string `json:"country" validate:"required,not_exists"`
}

func (request createPoliticalPartyRequest) Validate(ctx context.Context, handler politicalPartyHandler) error {
	_ = handler.v.Validator.RegisterTranslation("exists", handler.v.Translator, func(ut ut.Translator) error {
		return ut.Add("exists", "{0} already exist with country", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("exists", fe.Field())
		return t
	})
	handler.v.Validator.RegisterValidation("exists", func(fl validator.FieldLevel) bool {
		fieldName := fl.FieldName()
		fieldValue := fl.Field().String()
		if fieldName == "Slug" {
			value, _ := handler.service.Exists(ctx, map[string]interface{}{
				"slug":    fieldValue,
				"country": request.Country,
			}, nil)

			if value {
				return false
			}
			return true
		}
		if fieldName == "Name" {
			value, _ := handler.service.Exists(ctx, map[string]interface{}{
				"name":    fieldValue,
				"country": request.Country,
			}, nil)
			if value {
				return false
			}
			return true
		}
		return true
	})

	handler.v.Validator.RegisterValidation("not_exists", func(fl validator.FieldLevel) bool {
		fieldName := fl.FieldName()
		if fieldName == "Country" {
			value, _ := handler.countryService.IdExists(ctx, fl.Field().String())

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

type updatePoliticalPartyRequest struct {
	Name    string `json:"name" validate:"exists"`
	Slug    string `json:"slug" validate:"exists"`
	Country string `json:"country" validate:"not_exists"`
}

func (request updatePoliticalPartyRequest) Validate(ctx context.Context, handler politicalPartyHandler, params requestParams) error {
	_ = handler.v.Validator.RegisterTranslation("exists", handler.v.Translator, func(ut ut.Translator) error {
		return ut.Add("exists", "{0} already exist with country", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("exists", fe.Field())
		return t
	})
	handler.v.Validator.RegisterValidation("exists", func(fl validator.FieldLevel) bool {
		fieldName := fl.FieldName()
		fieldValue := fl.Field().String()
		if fieldName == "Slug" && fieldValue != "" {
			value, _ := handler.service.Exists(
				ctx,
				map[string]interface{}{
					"slug":    fieldValue,
					"country": request.Country,
				},
				map[string]interface{}{
					"_id": params.Id,
				},
			)

			if value {
				return false
			}
			return true
		}
		if fieldName == "Name" && fieldValue != "" {
			value, _ := handler.service.Exists(
				ctx,
				map[string]interface{}{
					"name":    fieldValue,
					"country": request.Country,
				},
				map[string]interface{}{
					"_id": params.Id,
				},
			)
			if value {
				return false
			}
			return true
		}
		return true
	})

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
		return true
	})

	err := handler.v.Validator.Struct(request)

	return err
}

type requestParams struct {
	Id string `uri:"id" validate:"required"`
}
