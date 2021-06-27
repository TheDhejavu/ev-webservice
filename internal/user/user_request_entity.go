package user

import (
	"context"

	validator "github.com/go-playground/validator/v10"
)

type createUserRequest struct {
	Username string `json:"username" validate:"required,min=1,exists"`
	Email    string `json:"email" validate:"required,email,exists"`
	Password string `json:"password" validate:"required,password"`
	Role     string `json:"role" validate:"required,oneof=user admin"`
	FullName string `json:"fullname" validate:"required,min=6"`
}

func (request *createUserRequest) Validate(ctx context.Context, handler userHandler) error {
	handler.v.Validator.RegisterValidation("exists", func(fl validator.FieldLevel) bool {
		fieldName := fl.FieldName()
		if fieldName == "Email" {
			value, _ := handler.service.IsEmailTaken(ctx, fl.Field().String(), nil)
			if value {
				return false
			}
			return true
		}
		if fieldName == "Username" {
			value, _ := handler.service.IsUsernameTaken(ctx, fl.Field().String(), nil)
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

type userRequestParams struct {
	Id string `uri:"id" validate:"required"`
}

type updateUserRequest struct {
	FullName string `json:"fullname" validate:"min=6"`
	Username string `json:"username" validate:"exists"`
	Email    string `json:"email" validate:"email,exists"`
	Password string `json:"password" validate:"password"`
	Role     string `json:"role"`
}

func (request *updateUserRequest) Validate(ctx context.Context, handler userHandler, params userRequestParams) error {
	handler.v.Validator.RegisterValidation("exists", func(fl validator.FieldLevel) bool {
		fieldName := fl.FieldName()
		fieldValue := fl.Field().String()
		if fieldName == "Email" && fieldValue != "" {
			value, _ := handler.service.IsEmailTaken(ctx, fl.Field().String(), params.Id)
			if value {
				return false
			}
			return true
		}
		if fieldName == "Username" && fieldValue != "" {
			value, _ := handler.service.IsUsernameTaken(ctx, fl.Field().String(), params.Id)
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
