package validator

import (
	"errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
)

type Validator struct {
	Validator  *validator.Validate
	Translator ut.Translator
}

func (v *Validator) Validate(i any) error {
	err := v.Validator.Struct(i)

	if err != nil {
		object, _ := err.(validator.ValidationErrors)
		for _, e := range object {
			log.Errorf("[Validate-1] %s: %s", e.Field(), e.Translate(v.Translator))

			return errors.New(e.Translate(v.Translator))

		}
	}

	return nil
}

func NewValidator() *Validator {
	en := en.New()
	uni := ut.New(en, en)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatalf("[NewValidator-1] Translator not found")
	}

	validate := validator.New()

	validate.RegisterTranslation("eqfield", trans, func(ut ut.Translator) error {
		return ut.Add("eqfield", "{0} and {1} do not match", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("eqfield", fe.Param(), fe.Field())

		return t
	})

	return &Validator{
		Validator:  validate,
		Translator: trans,
	}
}
