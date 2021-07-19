package utils

import (
	"log"
	"net/url"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomValidator struct {
	Validator  *validator.Validate
	Translator ut.Translator
}

var (
	uni *ut.UniversalTranslator
	v   *validator.Validate
)

func CustomValidators() *CustomValidator {
	translator := en.New()
	uni = ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	v = validator.New()

	en_translations.RegisterDefaultTranslations(v, trans)

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("exists", trans, func(ut ut.Translator) error {
		return ut.Add("exists", "{0} already exist", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("exists", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("not_exists", trans, func(ut ut.Translator) error {
		return ut.Add("not_exists", "{0} does not exist", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("not_exists", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("password", trans, func(ut ut.Translator) error {
		return ut.Add("password", "{0} is not strong enough", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("valid_country", trans, func(ut ut.Translator) error {
		return ut.Add("valid_country", "country is not valid", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("valid_country", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("url", trans, func(ut ut.Translator) error {
		return ut.Add("url", "url is not valid", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("url", fe.Field())
		return t
	})

	_ = v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > 6
	})

	_ = v.RegisterValidation("url", func(fl validator.FieldLevel) bool {
		fieldValue := fl.Field().String()
		_, err := url.ParseRequestURI(fieldValue)
		if err != nil {
			return false
		}
		return true
	})

	_ = v.RegisterValidation("valid_country", func(fl validator.FieldLevel) bool {
		fieldValue := fl.Field().String()
		_, err := primitive.ObjectIDFromHex(fieldValue)
		if err != nil {
			return false
		}
		return true
	})

	return &CustomValidator{v, trans}
}
