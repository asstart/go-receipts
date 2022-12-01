package declarativeusinggoplayground_test

import (
	"fmt"
	"testing"
	"time"

	vld "github.com/asstart/go-receipts/input_validation/declarative_using_goplayground"
	"github.com/go-playground/validator/v10"
)

func TestDefaultValidations(t *testing.T) {
	v := validator.New()
	tt := []struct {
		name  string
		obj   vld.ExampleStruct
		valid bool
	}{
		{"valid example", vld.ExampleStruct{ID: "12", Data: "data", CreatedAt: time.Now(), Count: 1}, true},
		{"invalid id length", vld.ExampleStruct{ID: "", Data: "data", CreatedAt: time.Now(), Count: 1}, false},
		{"invalid data length", vld.ExampleStruct{ID: "", Data: "123456", CreatedAt: time.Now(), Count: 1}, false},
		{"invalid count", vld.ExampleStruct{ID: "", Data: "data", CreatedAt: time.Now(), Count: 11}, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := v.Struct(tc.obj)
			if err == nil != tc.valid {
				t.Fatalf("got: %v", err)
			}
		})
	}
}

func TestCustomValidations(t *testing.T) {
	v := vld.CustomValidator()
	tt := []struct {
		name  string
		obj   vld.ExampleBase32Id
		valid bool
	}{
		{"valid base32 str", vld.ExampleBase32Id{ID: vld.GenBase32Id()}, true},
		{"invalid base32 str", vld.ExampleBase32Id{ID: "!?@"}, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := v.Struct(tc.obj)
			if err == nil != tc.valid {
				t.Fatalf("wanted: error_is_nil=%v, got: %v", tc.valid, err)
			}
		})
	}
}

func TestCustomTranslation(t *testing.T) {
	v := vld.CustomValidator()
	tr := vld.CustomTranslator(v)

	id := vld.ExampleBase32Id{ID: "!?@"}

	err := v.Struct(id)
	if err == nil {
		t.Fatal("err should be not nil")
	}

	errs, ok := err.(validator.ValidationErrors)
	fmt.Println(ok)
	if !ok {
		t.Fatal("error should be casted to validator.ValidationErrors")
	}

	if len(errs) != 1 {
		t.Fatalf("validation err should contains 1 elem")
	}
	expErr := "ID must be encoded as base32"
	actErr := errs[0].Translate(tr)
	if actErr != expErr {
		t.Fatalf("wanted: %v, got: %v", expErr, actErr)
	}
}
