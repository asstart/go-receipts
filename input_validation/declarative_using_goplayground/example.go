package declarativeusinggoplayground

import (
	"crypto/rand"
	"encoding/base32"
	"io"
	"log"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ExampleStruct struct {
	ID        string `validate:"len=2,required"`
	Data      string `validate:"max=4"`
	CreatedAt time.Time
	Count     int `validate:"max=10"`
}

func Base32IdValidator(fl validator.FieldLevel) bool {
	enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	_, err := enc.DecodeString(fl.Field().String())
	return err == nil
}

type ExampleBase32Id struct {
	ID string `validate:"baseid32"`
}

func CustomValidator() *validator.Validate {
	v := validator.New()

	err := v.RegisterValidation("baseid32", Base32IdValidator)
	if err != nil {
		log.Fatalf("error registrating validator: %v", err)
	}

	return v
}

func CustomTranslator(v *validator.Validate) ut.Translator {
	en := en.New()
	uni := ut.New(en, en)

	trans, _ := uni.GetTranslator("en")

	err := v.RegisterTranslation("baseid32", trans,
		func(tr ut.Translator) error {
			return tr.Add("baseid32", "{0} must be encoded as base32", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("baseid32", fe.Field())
			return t
		},
	)

	if err != nil {
		log.Fatalf("error registrating translation: %v", err)
	}

	return trans
}

func GenBase32Id() string {
	enc := base32.StdEncoding.WithPadding(base32.NoPadding)
	id := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, id)
	if err != nil {
		return ""
	}
	return enc.EncodeToString(id)
}
