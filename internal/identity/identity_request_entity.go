package identity

import (
	"context"
	"mime/multipart"

	validator "github.com/go-playground/validator/v10"
)

type createIdentityRequest struct {
	FirstName string `form:"first_name" validate:"required"`
	LastName  string `form:"last_name" validate:"required"`
	Origin    struct {
		Country string `form:"country" validate:"required,not_exists"`
		State   string `form:"state" validate:"required"`
		City    string `form:"city" validate:"required"`
		Address string `form:"address" validate:"required"`
	} `form:"origin" validate:"required,dive"`
	Residence struct {
		Country string `form:"country" validate:"required,not_exists"`
		State   string `form:"state" validate:"required"`
		City    string `form:"city" validate:"required"`
		Address string `form:"address" validate:"required"`
	} `form:"residence" validate:"required,dive"`
	Email            string                `form:"email" validate:"required,email"`
	Password         string                `form:"password" validate:"required"`
	BirthCertificate *multipart.FileHeader `form:"birth_certificate" bson:"birth_certificate" validate:"required"`
	NationalIdCard   *multipart.FileHeader `form:"national_id_card" bson:"national_id_card" validate:"required"`
	VoterCard        *multipart.FileHeader `form:"voter_card" bson:"voter_card" validate:"required"`
}

func (request createIdentityRequest) Validate(ctx context.Context, handler identityHandler) error {
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

type updateIdentityRequest struct {
	FirstName string `form:"first_name" bson:"first_name,omitempty"`
	LastName  string `form:"last_name" bson:"last_name,omitempty"`
	Origin    struct {
		Country string `form:"country" bson:"country,omitempty"`
		State   string `form:"state" bson:"state,omitempty"`
		City    string `form:"city" bson:"city,omitempty"`
		Address string `form:"address" bson:"address,omitempty"`
	} `form:"origin" bson:"origin,omitempty"`
	Residence *struct {
		Country string `form:"country" bson:"city,omitempty"`
		City    string `form:"city" bson:"city,omitempty"`
		State   string `form:"state" bson:"state,omitempty"`
		Address string `form:"address" bson:"address,omitempty"`
	} `form:"residence" bson:"residence,omitempty" validate:"required, dive"`
	Email            string                `form:"email" bson:"email"`
	Password         string                `form:"password" bson:"password"`
	BirthCertificate *multipart.FileHeader `form:"birth_certificate" bson:"birth_certificate"`
	NationalIdCard   *multipart.FileHeader `form:"national_id_card" bson:"national_id_card"`
	VoterCard        *multipart.FileHeader `form:"voter_card" bson:"voter_card"`
}

func (request updateIdentityRequest) Validate(ctx context.Context, handler identityHandler, params requestParams) error {

	return nil
}

type requestParams struct {
	Id string `uri:"id" validate:"required"`
}
