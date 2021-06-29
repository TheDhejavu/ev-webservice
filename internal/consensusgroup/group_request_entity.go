package consensusgroup

import (
	"context"

	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
)

type createGroupRequest struct {
	Name      string `json:"name" validate:"required,exists"`
	ServerUrl string `json:"server_url" validate:"required,url"`
	PublicKey string `json:"public_key" validate:"required"`
	Country   string `json:"country" validate:"valid_country,required,not_exists"`
}

func (request createGroupRequest) Validate(ctx context.Context, handler GroupHandler) error {
	_ = handler.v.Validator.RegisterTranslation("exists", handler.v.Translator, func(ut ut.Translator) error {
		return ut.Add("exists", "{0} already exist with country", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("exists", fe.Field())
		return t
	})
	handler.v.Validator.RegisterValidation("exists", func(fl validator.FieldLevel) bool {
		fieldName := fl.FieldName()
		fieldValue := fl.Field().String()
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

type updateGroupRequest struct {
	Name      string `json:"name" validate:"exists"`
	ServerUrl string `json:"server_url" validate:"url"`
	PublicKey string `json:"public_key"`
	Country   string `json:"country" validate:"valid_country,not_exists"`
}

func (request updateGroupRequest) Validate(ctx context.Context, handler GroupHandler, params requestParams) error {
	_ = handler.v.Validator.RegisterTranslation("exists", handler.v.Translator, func(ut ut.Translator) error {
		return ut.Add("exists", "{0} already exist with country", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("exists", fe.Field())
		return t
	})
	handler.v.Validator.RegisterValidation("exists", func(fl validator.FieldLevel) bool {
		fieldName := fl.FieldName()
		fieldValue := fl.Field().String()
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
